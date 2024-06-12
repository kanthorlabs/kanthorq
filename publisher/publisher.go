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
	pub.mu.Lock()
	defer pub.mu.Unlock()

	conn, err := pgx.Connect(ctx, pub.conf.ConnectionUri)
	if err != nil {
		return err
	}
	pub.conn = conn

	stream, err := q.Stream(ctx, pub.conn, &entities.Stream{Name: pub.conf.StreamName})
	if err != nil {
		return err
	}
	pub.stream = stream

	return nil
}

func (pub *publisher) Stop(ctx context.Context) error {
	pub.mu.Lock()
	defer pub.mu.Unlock()

	return pub.conn.Close(ctx)
}

func (pub *publisher) Send(ctx context.Context, events []*entities.StreamEvent) error {
	pub.mu.Lock()
	defer pub.mu.Unlock()

	ctx, span := telemetry.Tracer.Start(ctx, "publisher.Send")
	defer span.End()

	// wait for the transaction is done
	select {
	case <-ctx.Done():
		span.SetAttributes(attribute.String("ERROR.PUBLISHER.CONTEXT", ctx.Err().Error()))
		return ctx.Err()
	default:
		tx, err := pub.conn.Begin(ctx)
		if err != nil {
			return err
		}
		defer tx.Rollback(ctx)

		_, err = api.StreamEventPush(pub.stream, events).Do(ctx, tx)
		if err != nil {
			return err
		}

		return tx.Commit(ctx)
	}
}
