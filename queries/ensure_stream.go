package queries

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kanthorlabs/kanthorq/entities"
)

func EnsureStream(name string) func(ctx context.Context, conn *pgxpool.Pool) (*entities.Stream, error) {
	statement := `SELECT name, created_at, updated_at FROM stream_ensure(@stream_name);`
	args := pgx.NamedArgs{
		"stream_name": name,
	}

	return func(ctx context.Context, conn *pgxpool.Pool) (*entities.Stream, error) {
		var stream entities.Stream
		err := conn.
			QueryRow(ctx, statement, args).
			Scan(&stream.Name, &stream.CreatedAt, &stream.UpdatedAt)

		return &stream, err
	}
}
