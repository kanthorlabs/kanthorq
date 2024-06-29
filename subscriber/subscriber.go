package subscriber

import (
	"context"
	"errors"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/q"
	"github.com/kanthorlabs/kanthorq/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type subscriber struct {
	Conf     *Config
	Conn     *pgx.Conn
	Consumer *entities.Consumer
	Type     string

	errorc chan error
	mu     sync.Mutex
}

func (sub *subscriber) Start(ctx context.Context) error {
	if err := sub.Connect(ctx); err != nil {
		return err
	}

	sub.mu.Lock()
	defer sub.mu.Unlock()
	consumer, err := q.Consumer(ctx, sub.Conn, &entities.Consumer{
		StreamName: sub.Conf.StreamName,
		Topic:      sub.Conf.Topic,
		Name:       sub.Conf.ConsumerName,
	})
	if err != nil {
		return err
	}
	sub.Consumer = consumer

	return nil
}

func (sub *subscriber) Connect(ctx context.Context) error {
	sub.mu.Lock()
	defer sub.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if sub.Conn != nil && !sub.Conn.IsClosed() {
			return nil
		}

		conn, err := pgx.Connect(ctx, sub.Conf.ConnectionUri)
		if err != nil {
			return err
		}
		sub.Conn = conn

		return nil
	}
}

func (sub *subscriber) Stop(ctx context.Context) error {
	sub.mu.Lock()
	defer sub.mu.Unlock()

	// wait for errorc channel to be closed
	<-sub.Error()
	return sub.Conn.Close(ctx)
}

func (sub *subscriber) Error() <-chan error {
	return sub.errorc
}

func (sub *subscriber) Handle(ctx context.Context, handler SubscriberHandler, events []*entities.StreamEvent) {
	ctx, span := telemetry.Tracer().Start(ctx, "subscriber_consume_handle")
	span.SetAttributes(attribute.Int("event_count", len(events)))

	var cancelled []*entities.StreamEvent
	var retryable []*entities.StreamEvent
	var completed []*entities.StreamEvent

	// loop through each event to guarantee the order of events
	for _, event := range events {
		// continue our tracing for each event
		eventCtx := telemetry.
			Propagator().
			Extract(context.TODO(), telemetry.MapCarrier(event.Metadata))

		eventCtx, eventSpan := telemetry.Tracer().Start(eventCtx, "subscriber_consume_event", trace.WithSpanKind(trace.SpanKindConsumer))
		eventSpan.SetAttributes(attribute.String("stream_name", sub.Conf.StreamName))
		eventSpan.SetAttributes(attribute.String("consumer_name", sub.Conf.ConsumerName))
		eventSpan.SetAttributes(attribute.String("consumer_topic", sub.Conf.Topic))
		eventSpan.SetAttributes(attribute.String("event_id", event.EventId))

		state := handler(eventCtx, event)

		switch state {
		case entities.StateCancelled:
			cancelled = append(cancelled, event)
		case entities.StateCompleted:
			completed = append(completed, event)
		case entities.StateRetryable:
			retryable = append(retryable, event)
		}

		eventSpan.End()
	}
	span.SetAttributes(attribute.Int("event_tobe_cancel", len(cancelled)))
	span.SetAttributes(attribute.Int("event_tobe_complete", len(completed)))
	span.SetAttributes(attribute.Int("event_tobe_retry", len(retryable)))

	telemetry.MeterCounter("kanthorq_subscriber_consume_total")(
		int64(len(cancelled)),
		metric.WithAttributes(
			attribute.String("subscriber_type", sub.Type),
			attribute.String("stream_name", sub.Conf.StreamName),
			attribute.String("consumer_name", sub.Conf.ConsumerName),
			attribute.String("consumer_topic", sub.Conf.Topic),
			attribute.String("subscriber_status", entities.StateCancelled.String()),
		),
	)
	telemetry.MeterCounter("kanthorq_subscriber_consume_total")(
		int64(len(completed)),
		metric.WithAttributes(
			attribute.String("subscriber_type", sub.Type),
			attribute.String("stream_name", sub.Conf.StreamName),
			attribute.String("consumer_name", sub.Conf.ConsumerName),
			attribute.String("consumer_topic", sub.Conf.Topic),
			attribute.String("subscriber_status", entities.StateCompleted.String()),
		),
	)
	telemetry.MeterCounter("kanthorq_subscriber_consume_total")(
		int64(len(retryable)),
		metric.WithAttributes(
			attribute.String("subscriber_type", sub.Type),
			attribute.String("stream_name", sub.Conf.StreamName),
			attribute.String("consumer_name", sub.Conf.ConsumerName),
			attribute.String("consumer_topic", sub.Conf.Topic),
			attribute.String("subscriber_status", entities.StateRetryable.String()),
		),
	)

	tx, err := sub.Conn.Begin(ctx)
	if err != nil {
		sub.RecordError(span, err)
		return
	}

	if len(cancelled) > 0 {
		command := q.NewConsumerJobMarkCancelled(sub.Consumer, completed)
		if _, err := command.Do(ctx, tx); err != nil {
			sub.RecordError(span, err, tx.Rollback(ctx))
			return
		}
	}

	if len(completed) > 0 {
		command := q.NewConsumerJobMarkCompleted(sub.Consumer, completed)
		if _, err := command.Do(ctx, tx); err != nil {
			sub.RecordError(span, err, tx.Rollback(ctx))
			return
		}
	}

	if len(retryable) > 0 {
		command := q.NewConsumerJobMarkRetryable(sub.Consumer, retryable)
		if _, err := command.Do(ctx, tx); err != nil {
			sub.RecordError(span, err, tx.Rollback(ctx))
			return
		}
	}

	if err := tx.Commit(ctx); err != nil {
		sub.RecordError(span, err, tx.Rollback(ctx))
		return
	}
}

func (sub *subscriber) RecordError(span trace.Span, errs ...error) error {
	var merged error
	for _, err := range errs {
		if err == nil {
			continue
		}

		span.RecordError(err)
		merged = errors.Join(merged, err)
		sub.errorc <- err
	}
	return merged
}

func (sub *subscriber) ErrorHandle(handle func(error)) {
	for err := range sub.errorc {
		handle(err)
	}
}
