package core

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
)

type Doable[T any] interface {
	Do(ctx context.Context, tx pgx.Tx) (*T, error)
}

func Do[T any](ctx context.Context, req Doable[T], conn *pgx.Conn) (*T, error) {
	// there is no auto-rollback on context cancellation.
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	res, err := req.Do(ctx, tx)
	if err != nil {
		return nil, errors.Join(err, tx.Rollback(ctx))
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return res, nil
}

func DoWithCM[T any](ctx context.Context, req Doable[T], cm pgcm.ConnectionManager) (*T, error) {
	conn, err := cm.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer cm.Release(ctx, conn)

	return Do(ctx, req, conn)
}
