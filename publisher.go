package kanthorq

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kanthorlabs/kanthorq/api"
	"github.com/kanthorlabs/kanthorq/entities"
)

var _ Publisher = (*publisher)(nil)

func Pub(ctx context.Context, pool *pgxpool.Pool, name string) (Publisher, error) {
	stream, err := Stream(ctx, pool, name)
	if err != nil {
		return nil, err
	}

	return &publisher{stream: stream, pool: pool}, nil
}

type Publisher interface {
	Send(ctx context.Context, events []*entities.StreamEvent) error
}

type publisher struct {
	stream *entities.Stream
	pool   *pgxpool.Pool
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
