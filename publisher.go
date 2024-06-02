package kanthorq

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kanthorlabs/kanthorq/api"
	"github.com/kanthorlabs/kanthorq/entities"
)

var _ Publisher = (*publisher)(nil)

func Pub(ctx context.Context, pool *pgxpool.Pool, streamName string) (Publisher, error) {
	stream, err := Stream(ctx, pool, streamName)
	if err != nil {
		return nil, err
	}

	return &publisher{stream: stream, pool: pool}, nil
}

type Publisher interface {
	Send(ctx context.Context, events []*entities.StreamEvent) error
}

type publisher struct {
	pool *pgxpool.Pool

	stream *entities.Stream
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
