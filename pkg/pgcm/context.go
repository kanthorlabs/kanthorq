package pgcm

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type ctxkey string

var CONNECTION = ctxkey("pgcm.connection")

func WithConn(ctx context.Context, conn *pgx.Conn) context.Context {
	return context.WithValue(ctx, CONNECTION, conn)
}

func ConnFromCtx(ctx context.Context) (*pgx.Conn, bool) {
	conn, ok := ctx.Value(CONNECTION).(*pgx.Conn)
	return conn, ok
}
