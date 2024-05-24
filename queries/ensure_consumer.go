package queries

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kanthorlabs/kanthorq/entities"
)

func EnsureConsumer(name, stream, topic string) func(ctx context.Context, conn *pgxpool.Pool) (*entities.Consumer, error) {
	statement := `SELECT name, stream_name, topic, cursor, created_at, updated_at FROM consumer_ensure(@consumer_name, @stream_name, @topic);`
	args := pgx.NamedArgs{
		"consumer_name": name,
		"stream_name":   stream,
		"topic":         topic,
	}
	return func(ctx context.Context, conn *pgxpool.Pool) (*entities.Consumer, error) {
		var consumer entities.Consumer
		err := conn.
			QueryRow(ctx, statement, args).
			Scan(&consumer.Name, &consumer.StreamName, &consumer.Topic, &consumer.Cursor, &consumer.CreatedAt, &consumer.UpdatedAt)

		return &consumer, err
	}
}
