package queries

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
)

func EnsureStream(name string) func(ctx context.Context, tx pgx.Tx) (*entities.Stream, error) {
	statement := `
INSERT INTO kanthorq_stream(name)
VALUES(@stream_name)
ON CONFLICT(name) DO UPDATE 
SET updated_at = EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000
RETURNING name, created_at, updated_at;
	`

	creation := `
CREATE TABLE IF NOT EXISTS %s (
	topic VARCHAR(128) NOT NULL,
	event_id VARCHAR(64) NOT NULL,
	created_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000,
	PRIMARY KEY (topic, event_id)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_event_id ON %s USING btree("event_id");
`
	stable := pgx.Identifier{entities.CollectionStreamEvent(name)}.Sanitize()
	query := fmt.Sprintf(creation, stable, stable)

	return func(ctx context.Context, tx pgx.Tx) (*entities.Stream, error) {
		var stream entities.Stream
		err := tx.
			QueryRow(ctx, statement, pgx.NamedArgs{"stream_name": name}).
			Scan(&stream.Name, &stream.CreatedAt, &stream.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if _, err := tx.Exec(ctx, query); err != nil {
			return nil, err
		}

		return &stream, nil
	}
}
