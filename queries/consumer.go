package queries

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
)

func ConsumerPull(consumer *entities.Consumer, size int) func(ctx context.Context, tx pgx.Tx) (string, error) {
	ctable := pgx.Identifier{entities.CollectionConsumer}.Sanitize()
	lock := fmt.Sprintf(
		"SELECT cursor FROM %s WHERE name = @consumer_name FOR UPDATE SKIP LOCKED",
		ctable,
	)
	statement := `
WITH jobs AS (
	INSERT INTO %s (event_id, topic)
	SELECT event_id, topic
	FROM %s
	WHERE topic = @consumer_topic AND event_id > @consumer_cursor
	ORDER BY topic, event_id
	LIMIT @size
	ON CONFLICT(event_id) DO UPDATE
	SET updated_at = EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000
	RETURNING event_id
),
next_cursor AS (
	SELECT max(event_id) AS next_event FROM jobs
)
UPDATE kanthorq_consumer
SET cursor = next_event
FROM next_cursor
WHERE name = @consumer_name AND next_event IS NOT NULL
RETURNING cursor;	
`
	jtable := pgx.Identifier{entities.CollectionConsumerJob(consumer.Name)}.Sanitize()
	stable := pgx.Identifier{entities.CollectionStreamEvent(consumer.StreamName)}.Sanitize()
	query := fmt.Sprintf(statement, jtable, stable)

	return func(ctx context.Context, tx pgx.Tx) (string, error) {
		var cursor *string
		err := tx.
			QueryRow(ctx, lock, pgx.NamedArgs{"consumer_name": consumer.Name}).
			Scan(&cursor)
		if err != nil {
			return "", err
		}
		if cursor == nil {
			return "", errors.Join(fmt.Errorf("ERROR.CONSUMER.BUSY: %s", consumer.Name))
		}

		args := pgx.NamedArgs{
			"consumer_topic":  consumer.Topic,
			"consumer_cursor": *cursor,
			"size":            size,
			"consumer_name":   consumer.Name,
		}
		var next *string
		if err := tx.QueryRow(ctx, query, args).Scan(&next); err != nil {
			return "", err
		}
		// no more job
		if next == nil {
			return "", nil
		}

		// update cache data
		consumer.Cursor = *next
		return *next, nil
	}
}
