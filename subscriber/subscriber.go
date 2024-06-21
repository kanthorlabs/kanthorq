package subscriber

import (
	"context"
	"errors"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/api"
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

	failurec chan map[string]error
	errorc   chan error
	mu       sync.Mutex
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

	// wait for failurec and errorc to be closed
	<-sub.Error()
	<-sub.Failurec()
	return sub.Conn.Close(ctx)
}

func (sub *subscriber) Failurec() <-chan map[string]error {
	return sub.failurec
}

func (sub *subscriber) Error() <-chan error {
	return sub.errorc
}

func (sub *subscriber) Handle(ctx context.Context, handler SubscriberHandler, events []*entities.StreamEvent) {
	ctx, span := telemetry.Tracer().Start(ctx, "subscriber_consume_handle")
	span.SetAttributes(attribute.Int("event_count", len(events)))

	var failures = make(map[string]error)
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

		if err := handler(eventCtx, event); err != nil {
			failures[event.EventId] = err
			retryable = append(retryable, event)
			eventSpan.RecordError(err)
			eventSpan.End()
			continue
		}

		completed = append(completed, event)
		eventSpan.End()
	}
	span.SetAttributes(attribute.Int("event_tobe_retry", len(retryable)))
	span.SetAttributes(attribute.Int("event_tobe_complete", len(completed)))

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

	// no error reports, mark jobs as completed
	if len(completed) > 0 {
		command := api.NewConsumerJobMarkComplete(sub.Consumer, completed)
		if _, err := command.Do(ctx, tx); err != nil {
			sub.RecordError(span, err, tx.Rollback(ctx))
			return
		}
	}

	if len(retryable) > 0 {
		// error reports, mark jobs as retryable
		command := api.NewConsumerJobMarkRetry(sub.Consumer, retryable)
		if _, err := command.Do(ctx, tx); err != nil {
			sub.RecordError(span, err, tx.Rollback(ctx))
			return
		}
	}

	if err := tx.Commit(ctx); err != nil {
		sub.RecordError(span, err, tx.Rollback(ctx))
		return
	}
	sub.failurec <- failures
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

func (sub *subscriber) FailureHandle(handle func(string, error)) {
	for failures := range sub.failurec {
		for eventId, err := range failures {
			handle(eventId, err)
		}
	}
}
