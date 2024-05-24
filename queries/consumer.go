package queries

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kanthorlabs/kanthorq/entities"
)

func ConsumerPull(name string, size int) func(ctx context.Context, conn *pgxpool.Pool) (*entities.ConsumerCursor, error) {
	statement := `SELECT consumer_name, cursor_current, cursor_next FROM kanthorq_consumer_pull(@consumer_name, CAST(@size as SMALLINT));`
	args := pgx.NamedArgs{
		"consumer_name": name,
		"size":          size,
	}

	return func(ctx context.Context, conn *pgxpool.Pool) (*entities.ConsumerCursor, error) {
		var cursor entities.ConsumerCursor
		err := conn.
			QueryRow(ctx, statement, args).
			Scan(&cursor.Name, &cursor.Current, &cursor.Next)

		return &cursor, err
	}
}
