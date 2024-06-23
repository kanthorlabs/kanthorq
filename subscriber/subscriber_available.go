package subscriber

import (
	"context"
	"fmt"
	"time"

	"github.com/kanthorlabs/kanthorq/api"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var _ Subscriber = (*available)(nil)

func NewAvailable(conf *Config) Subscriber {
	return &available{
		subscriber: &subscriber{
			Conf:     conf,
			failurec: make(chan map[string]error, 100),
			errorc:   make(chan error, 100),
			Type:     fmt.Sprintf("subscribe_%s", entities.StateAvailable.String()),
		},
	}
}

type available struct {
	*subscriber
}

func (sub *available) Pull(ctx context.Context, options ...SubscribeOption) ([]*entities.StreamEvent, error) {
	opts := NewSubscribeOption(options...)

	ctx, span := telemetry.Tracer().Start(ctx, "subscriber_pull", trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()
	span.SetAttributes(attribute.String("subscribe_type", sub.Type))
	span.SetAttributes(attribute.String("stream_name", sub.Conf.StreamName))
	span.SetAttributes(attribute.String("topic", sub.Conf.Topic))
	span.SetAttributes(attribute.String("consumer_name", sub.Conf.ConsumerName))
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
	tx, err := sub.Conn.Begin(ctx)
	if err != nil {
		return nil, sub.RecordError(span, err)
	}

	c, err := api.NewConsumerPull(sub.Consumer, opts.Size).Do(ctx, tx)
	if err != nil {
		return nil, sub.RecordError(span, err, tx.Rollback(ctx))
	}

	// no more events to pull so we cannot set next cursor
	if c.NextCursor == "" {
		span.SetAttributes(attribute.Bool("err_no_event_in_stream", true))
		span.SetAttributes(attribute.Bool("err_next_cursor_empty", true))
		return nil, sub.RecordError(span, tx.Rollback(ctx))
	}

	j, err := api.NewConsumerJobPullAvailable(sub.Consumer, opts.Size, opts.VisibilityTimeout).Do(ctx, tx)
	if err != nil {
		return nil, sub.RecordError(span, err, tx.Rollback(ctx))
	}

	span.SetAttributes(attribute.Int("event_count", len(j.Events)))
	// no event was found
	if len(j.Events) == 0 {
		span.SetAttributes(attribute.Bool("err_no_event_in_consumer", true))
		return nil, sub.RecordError(span, tx.Rollback(ctx))
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, sub.RecordError(span, err)
	}

	telemetry.MeterCounter("kanthorq_subscriber_pull_total")(
		int64(len(j.Events)),
		metric.WithAttributes(
			attribute.String("subscriber_type", sub.Type),
			attribute.String("stream_name", sub.Conf.StreamName),
			attribute.String("consumer_name", sub.Conf.ConsumerName),
			attribute.String("consumer_topic", sub.Conf.Topic),
		),
	)
	return j.Events, nil
}

func (sub *available) Consume(ctx context.Context, handler SubscriberHandler, options ...SubscribeOption) {
	opts := NewSubscribeOption(options...)

	// start error handler
	go sub.ErrorHandle(opts.OnError)
	// start failure handler
	go sub.FailureHandle(opts.OnFailure)

	for {
		// cctx will be use for consume logic
		cctx, cancel := context.WithTimeout(context.Background(), opts.WaitingTime)
		cctx, span := telemetry.Tracer().Start(ctx, "subscriber_consume", trace.WithSpanKind(trace.SpanKindConsumer))
		span.SetAttributes(attribute.String("subscribe_type", sub.Type))
		span.SetAttributes(attribute.String("stream_name", sub.Conf.StreamName))
		span.SetAttributes(attribute.String("consumer_topic", sub.Conf.Topic))
		span.SetAttributes(attribute.String("consumer_name", sub.Conf.ConsumerName))
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
			if err := sub.Connect(cctx); err != nil {
				close(sub.errorc)
				close(sub.failurec)
				cancel()
				span.End()
				// if we still can't connect, throw it
				panic(err)
			}

			events, err := sub.Pull(ctx, options...)
			if err != nil {
				sub.RecordError(span, err)
				cancel()
				span.End()
				time.Sleep(opts.WaitingTime)
				return
			}

			sub.Handle(cctx, handler, events)
			cancel()
			span.End()
			time.Sleep(opts.WaitingTime)
		}
	}
}
