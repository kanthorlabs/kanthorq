package queries

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
)

func EnsureConsumer(name, streamName, topic string) func(ctx context.Context, tx pgx.Tx) (*entities.Consumer, error) {
	statement := `
INSERT INTO kanthorq_consumer(name, stream_name, topic, cursor)
VALUES(@consumer_name, @stream_name, @topic, '')
ON CONFLICT(name) DO UPDATE 
SET updated_at = EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000
RETURNING name, stream_name, topic, cursor, created_at, updated_at;
`

	creation := `
CREATE TABLE IF NOT EXISTS %s (
	event_id VARCHAR(64) NOT NULL,
	topic VARCHAR(128) NOT NULL,
	state SMALLINT NOT NULL DEFAULT 0,
	schedule_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000,
	attempt_count SMALLINT NOT NULL DEFAULT 0,
	attempted_at BIGINT NOT NULL DEFAULT 0,
	created_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000,
	updated_at BIGINT NOT NULL DEFAULT 0,
	PRIMARY KEY (event_id)
);
`
	jtable := pgx.Identifier{entities.CollectionConsumerJob(name)}.Sanitize()
	query := fmt.Sprintf(creation, jtable)

	return func(ctx context.Context, tx pgx.Tx) (*entities.Consumer, error) {
		var consumer entities.Consumer
		var args = pgx.NamedArgs{
			"consumer_name": name,
			"stream_name":   streamName,
			"topic":         topic,
		}
		err := tx.
			QueryRow(ctx, statement, args).
			Scan(
				&consumer.Name,
				&consumer.StreamName,
				&consumer.Topic,
				&consumer.Cursor,
				&consumer.CreatedAt,
				&consumer.UpdatedAt,
			)
		if err != nil {
			return nil, err
		}

		if _, err := tx.Exec(ctx, query); err != nil {
			return nil, err
		}

		return &consumer, nil
	}
}
