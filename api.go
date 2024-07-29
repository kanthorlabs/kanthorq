package kanthorq

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
)

type Doable[T any] interface {
	Do(ctx context.Context, tx pgx.Tx) (*T, error)
}

func Do[T any](ctx context.Context, req Doable[T], conn *pgx.Conn) (*T, error) {
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	res, err := req.Do(ctx, tx)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return res, nil
}

func DoWithCM[T any](ctx context.Context, req Doable[T], cm pgcm.ConnectionManager) (*T, error) {
	conn, err := cm.Connection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	return Do(ctx, req, conn.Raw())
}
