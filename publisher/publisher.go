package publisher

import (
	"context"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/q"
	"github.com/kanthorlabs/kanthorq/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var _ Publisher = (*publisher)(nil)

func New(conf *Config) Publisher {
	return &publisher{conf: conf}
}

type publisher struct {
	conf *Config
	mu   sync.Mutex

	conn   *pgx.Conn
	stream *entities.Stream
}

func (pub *publisher) Start(ctx context.Context) error {
	if err := pub.connect(ctx); err != nil {
		return err
	}

	pub.mu.Lock()
	defer pub.mu.Unlock()
	stream, err := q.NewStream(ctx, pub.conn, &entities.Stream{Name: pub.conf.StreamName})
	if err != nil {
		return err
	}
	pub.stream = stream

	return nil
}

func (pub *publisher) connect(ctx context.Context) error {
	pub.mu.Lock()
	defer pub.mu.Unlock()

	// @TODO: test what will happen if pgbouncer terminate a connection
	if pub.conn != nil && !pub.conn.IsClosed() {
		return nil
	}

	conn, err := pgx.Connect(ctx, pub.conf.ConnectionUri)
	if err != nil {
		return err
	}
	pub.conn = conn

	return nil
}
func (pub *publisher) Stop(ctx context.Context) error {
	pub.mu.Lock()
	defer pub.mu.Unlock()

	return pub.conn.Close(ctx)
}

func (pub *publisher) Send(ctx context.Context, events []*entities.StreamEvent) error {
	ctx, span := telemetry.Tracer().Start(ctx, "publisher_send", trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	// @TODO: good time to introduce chain of responsibility pattern -> middleware
	// - add circuit breaker
	// - add retry logic
	// - add throttle logic
	start := time.Now()
	defer func() {
		// elapsed.Seconds may return 0 if the duration is too small
		// so devide milliseconds by 1000 to get seconds is better
		duration := float64(time.Since(start).Milliseconds()) / 1000
		telemetry.MeterHistogram("kanthorq_publisher_send_duration_seconds")(
			duration,
			metric.WithAttributes(attribute.String("stream_name", pub.conf.StreamName)),
		)
	}()

	// wait for the transaction is done
	select {
	case <-ctx.Done():
		span.RecordError(ctx.Err())
		return ctx.Err()
	default:
		if err := pub.connect(ctx); err != nil {
			span.RecordError(err)
			return err
		}

		// TODO: add retry logic
		tx, err := pub.conn.Begin(ctx)
		if err != nil {
			span.SetAttributes(attribute.Bool("ErrTxBegin", true))
			span.RecordError(err)
			return err
		}

		_, err = q.StreamEventPush(pub.stream, events).Do(ctx, tx)
		if err != nil {
			span.SetAttributes(attribute.Bool("ErrApiStreamEventPush", true))
			span.RecordError(err)

			if err := tx.Rollback(ctx); err != nil {
				span.SetAttributes(attribute.Bool("ErrTxRollback", true))
				span.RecordError(err)
			}
			return err
		}

		if err := tx.Commit(ctx); err != nil {
			span.SetAttributes(attribute.Bool("ErrTxCommit", true))
			span.RecordError(err)

			return err
		}

		for _, event := range events {
			telemetry.MeterCounter("kanthorq_publisher_send_total")(
				1,
				metric.WithAttributes(
					attribute.String("stream_name", pub.conf.StreamName),
					attribute.String("consumer_topic", event.Topic),
				),
			)
		}
		return nil
	}
}
