package publisher

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/api"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/q"
	"github.com/kanthorlabs/kanthorq/telemetry"
	"go.opentelemetry.io/otel/attribute"
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
	stream, err := q.Stream(ctx, pub.conn, &entities.Stream{Name: pub.conf.StreamName})
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
	ctx, span := telemetry.Tracer.Start(ctx, "publisher.Send", trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

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

		_, err = api.StreamEventPush(pub.stream, events).Do(ctx, tx)
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

		return nil
	}
}
