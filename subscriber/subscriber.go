package subscriber

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/api"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/q"
	"github.com/kanthorlabs/kanthorq/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var _ Subscriber = (*subscriber)(nil)

func New(conf *Config) Subscriber {
	return &subscriber{
		conf: conf,
		// don't use unbuffer channel because it will block .Consume method
		// because when you send a value on an unbuffered channel,
		// the sending goroutine is blocked until another goroutine receives the value from the channel
		failurec: make(chan map[string]error, 1),
		errorc:   make(chan error, 1),
	}
}

type subscriber struct {
	conf     *Config
	failurec chan map[string]error
	errorc   chan error
	mu       sync.Mutex

	conn     *pgx.Conn
	consumer *entities.Consumer
}

func (sub *subscriber) Start(ctx context.Context) error {
	if err := sub.connect(ctx); err != nil {
		return err
	}

	sub.mu.Lock()
	defer sub.mu.Unlock()
	consumer, err := q.Consumer(ctx, sub.conn, &entities.Consumer{
		StreamName: sub.conf.StreamName,
		Topic:      sub.conf.Topic,
		Name:       sub.conf.ConsumerName,
	})
	if err != nil {
		return err
	}
	sub.consumer = consumer

	return nil
}

func (sub *subscriber) connect(ctx context.Context) error {
	sub.mu.Lock()
	defer sub.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if sub.conn != nil && !sub.conn.IsClosed() {
			return nil
		}

		conn, err := pgx.Connect(ctx, sub.conf.ConnectionUri)
		if err != nil {
			return err
		}
		sub.conn = conn

		return nil
	}
}

func (sub *subscriber) Stop(ctx context.Context) error {
	sub.mu.Lock()
	defer sub.mu.Unlock()

	// wait for failurec and errorc to be closed
	<-sub.Error()
	<-sub.Failurec()
	return sub.conn.Close(ctx)
}

func (sub *subscriber) Pull(ctx context.Context, options ...SubscribeOption) ([]*entities.StreamEvent, error) {
	opts := NewSubscribeOption(options...)

	ctx, span := telemetry.Tracer.Start(ctx, "subscriber.Pull", trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()
	span.SetAttributes(attribute.String("stream_name", sub.conf.StreamName))
	span.SetAttributes(attribute.String("topic", sub.conf.Topic))
	span.SetAttributes(attribute.String("consumer_name", sub.conf.ConsumerName))
	span.SetAttributes(attribute.Int("size", opts.Size))
	span.SetAttributes(attribute.Int64("timeout", opts.Timeout.Milliseconds()))
	span.SetAttributes(attribute.Int64("visibility_timeout", opts.VisibilityTimeout.Milliseconds()))
	span.SetAttributes(attribute.Int64("waiting_time", opts.WaitingTime.Milliseconds()))

	// if the parent context (from .Consume for example) is timeout
	// then this context will be timeout as well
	ctx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()

	// both .Begin and .Rollback will teriminate the underlying connection
	// if the underlying connection is closed or context timeout
	tx, err := sub.conn.Begin(ctx)
	if err != nil {
		return nil, sub.error(span, err)
	}

	c, err := api.NewConsumerPull(sub.consumer, opts.Size).Do(ctx, tx)
	// no new events to pull
	// we catch ErrNoRows as success case
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		span.SetAttributes(attribute.Bool("api.ConsumerPull/ErrNoRows", true))
		return nil, sub.error(span, tx.Rollback(ctx))
	}
	if err != nil {
		return nil, sub.error(span, err, tx.Rollback(ctx))
	}

	// no more events to pull
	if c.NextCursor == "" {
		span.SetAttributes(attribute.Bool("api.ConsumerPull/ErrNoNextCursor", true))
		return nil, sub.error(span, tx.Rollback(ctx))
	}

	j, err := api.NewConsumerJobPull(sub.consumer, opts.Size, opts.VisibilityTimeout).Do(ctx, tx)
	// no new events to pull
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		span.SetAttributes(attribute.Bool("api.ConsumerJobPull/ErrNoRows", true))
		return nil, sub.error(span, tx.Rollback(ctx))
	}
	if err != nil {
		return nil, sub.error(span, err, tx.Rollback(ctx))
	}

	span.SetAttributes(attribute.Int("api.ConsumerJobPull/Events", len(j.Events)))
	// no event was found
	if len(j.Events) == 0 {
		return nil, sub.error(span, tx.Rollback(ctx))
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, sub.error(span, err)
	}

	return j.Events, nil
}

func (sub *subscriber) Consume(ctx context.Context, handler SubscriberHandler, options ...SubscribeOption) {
	opts := NewSubscribeOption(options...)

	// start error handler
	go sub.errorh(opts.OnError)
	// start failure handler
	go sub.failureh(opts.OnFailure)

	for {
		// cctx will be use for consume logic
		cctx, cancel := context.WithTimeout(context.Background(), opts.WaitingTime)
		cctx, span := telemetry.Tracer.Start(ctx, "subscriber.Consume", trace.WithSpanKind(trace.SpanKindConsumer))
		span.SetAttributes(attribute.String("stream_name", sub.conf.StreamName))
		span.SetAttributes(attribute.String("topic", sub.conf.Topic))
		span.SetAttributes(attribute.String("consumer_name", sub.conf.ConsumerName))
		span.SetAttributes(attribute.Int("size", opts.Size))
		span.SetAttributes(attribute.Int64("timeout", opts.Timeout.Milliseconds()))
		span.SetAttributes(attribute.Int64("visibility_timeout", opts.VisibilityTimeout.Milliseconds()))
		span.SetAttributes(attribute.Int64("waiting_time", opts.WaitingTime.Milliseconds()))

		select {
		// ctx is used for control flow
		case <-ctx.Done():
			close(sub.errorc)
			close(sub.failurec)
			cancel()
			span.RecordError(ctx.Err())
			span.End()
			return
		default:
			// ctx is used for control flow
			if ctx.Err() != nil {
				close(sub.errorc)
				close(sub.failurec)
				cancel()
				span.RecordError(ctx.Err())
				span.End()
				return
			}

			// both .Begin and .Rollback will teriminate the underlying connection
			// if the underlying connection is closed or context timeout
			// so we need an helper to check our connection status before start consuming
			if err := sub.connect(cctx); err != nil {
				close(sub.errorc)
				close(sub.failurec)
				cancel()
				span.End()
				// if we still can't connect, throw it
				panic(err)
			}

			sub.consume(cctx, handler, options)
			cancel()
			span.End()
			time.Sleep(opts.WaitingTime)
		}
	}
}

func (sub *subscriber) consume(ctx context.Context, handler SubscriberHandler, options []SubscribeOption) {
	ctx, span := telemetry.Tracer.Start(ctx, "subscriber.Consume/consume", trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()

	events, err := sub.Pull(ctx, options...)
	if err != nil {
		sub.error(span, err)
		return
	}
	if len(events) == 0 {
		span.SetAttributes(attribute.Bool("subscriber.Consume/consume/NoEvents", true))
		return
	}

	var failures = make(map[string]error)
	var retryable []string
	var completed []string

	_, handlerSpan := telemetry.Tracer.Start(ctx, "subscriber.Consume/consume/handler")
	// loop through each event to guarantee the order of events
	for _, event := range events {
		// continue our tracing for each event
		continuosCtx := otel.
			GetTextMapPropagator().
			Extract(ctx, telemetry.MapCarrier(event.Metadata))

		eventCtx, eventSpan := telemetry.Tracer.Start(continuosCtx, "subscriber.Consume/consume/event", trace.WithSpanKind(trace.SpanKindConsumer))
		eventSpan.SetAttributes(attribute.String("event_id", event.EventId))

		if err := handler(eventCtx, event); err != nil {
			failures[event.EventId] = err
			retryable = append(retryable, event.EventId)
			eventSpan.RecordError(err)
			eventSpan.End()
			continue
		}

		completed = append(completed, event.EventId)
		eventSpan.End()
	}
	handlerSpan.End()

	tx, err := sub.conn.Begin(ctx)
	if err != nil {
		sub.error(span, err)
		return
	}

	// no error reports, mark jobs as completed
	if len(completed) > 0 {
		command := api.NewConsumerJobMarkComplete(sub.consumer, completed)
		if _, err := command.Do(ctx, tx); err != nil {
			sub.error(span, err, tx.Rollback(ctx))
			return
		}
	}

	if len(retryable) > 0 {
		// error reports, mark jobs as retryable
		command := api.NewConsumerJobMarkRetry(sub.consumer, retryable)
		if _, err := command.Do(ctx, tx); err != nil {
			sub.error(span, err, tx.Rollback(ctx))
			return
		}
	}

	if err := tx.Commit(ctx); err != nil {
		sub.error(span, err, tx.Rollback(ctx))
		return
	}
	sub.failurec <- failures
}

func (sub *subscriber) Failurec() <-chan map[string]error {
	return sub.failurec
}

func (sub *subscriber) Error() <-chan error {
	return sub.errorc
}

func (sub *subscriber) error(span trace.Span, errs ...error) error {
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

func (sub *subscriber) errorh(handle func(error)) {
	for err := range sub.errorc {
		handle(err)
	}
}

func (sub *subscriber) failureh(handle func(string, error)) {
	for failures := range sub.failurec {
		for eventId, err := range failures {
			handle(eventId, err)
		}
	}
}
