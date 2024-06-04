package publisher

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kanthorlabs/kanthorq/api"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/q"
)

var _ Publisher = (*publisher)(nil)

func New(conf *Config) Publisher {
	return &publisher{conf: conf}
}

type publisher struct {
	conf *Config

	pool   *pgxpool.Pool
	stream *entities.Stream
}

func (pub *publisher) Start(ctx context.Context) error {
	pool, err := pgxpool.New(ctx, pub.conf.ConnectionUri)
	if err != nil {
		return err
	}
	pub.pool = pool

	stream, err := q.Stream(ctx, pub.pool, &entities.Stream{Name: pub.conf.StreamName})
	if err != nil {
		return err
	}
	pub.stream = stream

	return nil
}

func (pub *publisher) Stop(ctx context.Context) error {
	pub.pool.Close()
	return nil
}

func (pub *publisher) Send(ctx context.Context, events []*entities.StreamEvent) error {
	tx, err := pub.pool.Begin(ctx)
	if err != nil {
		return err
	}

	_, err = api.StreamEventPush(pub.stream, events).Do(ctx, tx)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
