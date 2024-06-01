package kanthorq

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kanthorlabs/kanthorq/api"
	"github.com/kanthorlabs/kanthorq/entities"
)

func Stream(ctx context.Context, pool *pgxpool.Pool, name string) (*entities.Stream, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	r, err := api.StreamEnsure(name).Do(ctx, tx)
	if err != nil {
		return nil, err
	}

	return r.Stream, tx.Commit(ctx)
}
