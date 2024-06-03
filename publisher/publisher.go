package publisher

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kanthorlabs/kanthorq"
	"github.com/kanthorlabs/kanthorq/api"
	"github.com/kanthorlabs/kanthorq/entities"
)

var _ Publisher = (*publisher)(nil)

func New(conf *Config, pool *pgxpool.Pool) Publisher {
	return &publisher{conf: conf, pool: pool}
}

type publisher struct {
	conf *Config
	pool *pgxpool.Pool

	stream *entities.Stream
}

func (sub *publisher) Start(ctx context.Context) error {
	stream, err := kanthorq.Stream(ctx, sub.pool, &entities.Stream{
		Name: sub.conf.StreamName,
	})
	if err != nil {
		return err
	}

	sub.stream = stream
	return nil
}

func (sub *publisher) Stop(ctx context.Context) error {
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
