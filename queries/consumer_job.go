package queries

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
)

func ConsumerJobPull(consumer *entities.Consumer, size int) func(ctx context.Context, tx pgx.Tx) ([]*entities.StreamEvent, error) {
	// using tuple condition to force postgre uses our multiple column primary key
	// you also need to escape the character % in the statement WHERE (topic, event_id) IN (%%s) to rewrite it later
	statement := `
SELECT topic, event_id, created_at FROM %s AS stream
WHERE (topic, event_id) IN (%%s)
ORDER BY topic, event_id;
`
	stable := pgx.Identifier{entities.CollectionStreamEvent(consumer.StreamName)}.Sanitize()
	query := fmt.Sprintf(statement, stable)

	return func(ctx context.Context, tx pgx.Tx) ([]*entities.StreamEvent, error) {
		pks, err := ConsumerJobChangeState(consumer, size, entities.StateAvailable, entities.StateRunning)(ctx, tx)
		if err != nil {
			return nil, err
		}
		if len(pks) == 0 {
			return []*entities.StreamEvent{}, nil
		}

		var in = make([]string, len(pks))
		for i, pk := range pks {
			in[i] = fmt.Sprintf("('%s', '%s')", pk.Topic, pk.EventId)
		}

		// primary key pairs are safe to pass them directly to the query
		rewrite := fmt.Sprintf(query, strings.Join(in, ","))
		rows, err := tx.Query(ctx, rewrite)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		events := make([]*entities.StreamEvent, 0)
		for rows.Next() {
			var event entities.StreamEvent
			if err := rows.Scan(&event.Topic, &event.EventId, &event.CreatedAt); err != nil {
				return nil, err
			}
			events = append(events, &event)
		}

		return events, rows.Err()
	}
}

func ConsumerJobChangeState(consumer *entities.Consumer, size int, from, to entities.JobState) func(ctx context.Context, tx pgx.Tx) ([]*entities.EventPk, error) {
	statement := `
WITH locked_jobs AS (
	SELECT event_id
	FROM %s AS l_jobs
	WHERE l_jobs.state = @from_state AND l_jobs.schedule_at < @attempt_at
	ORDER BY l_jobs.event_id ASC
	LIMIT @size
	FOR UPDATE SKIP LOCKED
)
UPDATE %s AS u_jobs
SET
	state = @to_state,
	attempt_count = attempt_count + 1,
	attempted_at = @attempt_at,
	schedule_at = @next_schedule_at
FROM locked_jobs
WHERE u_jobs.event_id = locked_jobs.event_id
RETURNING u_jobs.topic, u_jobs.event_id
`
	jtable := pgx.Identifier{entities.CollectionConsumerJob(consumer.Name)}.Sanitize()
	query := fmt.Sprintf(statement, jtable, jtable)

	return func(ctx context.Context, tx pgx.Tx) ([]*entities.EventPk, error) {
		args := pgx.NamedArgs{
			"from_state":       from,
			"attempt_at":       time.Now().UTC().UnixMilli(),
			"to_state":         to,
			"size":             size,
			"next_schedule_at": time.Now().Add(time.Second * 3600).UTC().UnixMilli(),
		}
		rows, err := tx.Query(ctx, query, args)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var pks []*entities.EventPk
		for rows.Next() {
			var pk entities.EventPk
			err = rows.Scan(&pk.Topic, &pk.EventId)
			if err != nil {
				return nil, err
			}
			pks = append(pks, &pk)
		}

		return pks, rows.Err()
	}
}
