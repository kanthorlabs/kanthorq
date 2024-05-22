package core

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func QueryStreamEnsure(name string) func(ctx context.Context, conn *pgxpool.Pool) (*Stream, error) {
	statement := `SELECT name, created_at, updated_at FROM stream_ensure(@stream_name);`
	args := pgx.NamedArgs{
		"stream_name": name,
	}

	return func(ctx context.Context, conn *pgxpool.Pool) (*Stream, error) {
		var stream Stream
		err := conn.
			QueryRow(ctx, statement, args).
			Scan(&stream.Name, &stream.CreatedAt, &stream.UpdatedAt)

		return &stream, err
	}
}

func QueryConsumerEnsure(name, stream, topic string) func(ctx context.Context, conn *pgxpool.Pool) (*Consumer, error) {
	statement := `SELECT name, stream_name, topic, cursor, created_at, updated_at FROM consumer_ensure(@consumer_name, @stream_name, @topic);`
	args := pgx.NamedArgs{
		"consumer_name": name,
		"stream_name":   stream,
		"topic":         topic,
	}
	return func(ctx context.Context, conn *pgxpool.Pool) (*Consumer, error) {
		var consumer Consumer
		err := conn.
			QueryRow(ctx, statement, args).
			Scan(&consumer.Name, &consumer.StreamName, &consumer.Topic, &consumer.Cursor, &consumer.CreatedAt, &consumer.UpdatedAt)

		return &consumer, err
	}
}

func QueryConsumerPull(name string, size int) func(ctx context.Context, conn *pgxpool.Pool) (*ConsumerCursor, error) {
	statement := `SELECT consumer_name, cursor_current, cursor_next FROM kanthorq_consumer_pull(@consumer_name, CAST(@size as SMALLINT));`
	args := pgx.NamedArgs{
		"consumer_name": name,
		"size":          size,
	}

	return func(ctx context.Context, conn *pgxpool.Pool) (*ConsumerCursor, error) {
		var cursor ConsumerCursor
		err := conn.
			QueryRow(ctx, statement, args).
			Scan(&cursor.Name, &cursor.Current, &cursor.Next)

		return &cursor, err
	}
}
