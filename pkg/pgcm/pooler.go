package pgcm

import (
	"context"

	"github.com/jackc/pgx/v5"
)

var _ ConnectionManager = (*pooler)(nil)

// NewPooler initializes a ConnectionManager that connect to an external PG pooler
// PGBouncer for instance
// connections should be handled by the pooler instead of our client code
// so everytime we finish with a connection, we should return it to the pooler by closing it
func NewPooler(uri string) ConnectionManager {
	return &pooler{uri: uri}
}

type pooler struct {
	uri string
}

func (cm *pooler) Start(ctx context.Context) error {
	return nil
}

func (cm *pooler) Stop(ctx context.Context) error {
	return nil
}

func (cm *pooler) Acquire(ctx context.Context) (*pgx.Conn, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	conn, err := pgx.Connect(ctx, cm.uri)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (cm *pooler) Release(ctx context.Context, conn *pgx.Conn) error {
	return conn.Close(ctx)
}
