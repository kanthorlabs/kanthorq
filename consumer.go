package kanthorq

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kanthorlabs/kanthorq/api"
	"github.com/kanthorlabs/kanthorq/entities"
)

func Consumer(ctx context.Context, pool *pgxpool.Pool, consumer *entities.Consumer) (*entities.Consumer, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	s, err := api.StreamEnsure(consumer.StreamName).Do(ctx, tx)
	if err != nil {
		return nil, err
	}

	r, err := api.ConsumerEnsure(s.Stream, consumer.Name, consumer.Topic).Do(ctx, tx)
	if err != nil {
		return nil, err
	}

	return r.Consumer, tx.Commit(ctx)
}
