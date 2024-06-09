package api

import (
	"context"
	_ "embed"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
)

func ConsumerJobPull(consumer *entities.Consumer, size int, vt time.Duration) *ConsumerJobPullReq {
	return &ConsumerJobPullReq{
		Consumer:          consumer,
		Size:              size,
		VisibilityTimeout: vt,
	}
}

//go:embed consumer_job_pull.sql
var ConsumerJobPullSQL string

type ConsumerJobPullReq struct {
	Consumer          *entities.Consumer
	Size              int
	VisibilityTimeout time.Duration
}

type ConsumerJobPullRes struct {
	Events map[string]*entities.StreamEvent
}

func (req *ConsumerJobPullReq) Do(ctx context.Context, tx pgx.Tx) (*ConsumerJobPullRes, error) {
	command := ConsumerJobStateChange(
		req.Consumer,
		req.Size,
		entities.StateAvailable,
		entities.StateRunning,
		req.VisibilityTimeout,
	)
	changes, err := command.Do(ctx, tx)
	if err != nil {
		return nil, err
	}

	res := &ConsumerJobPullRes{Events: make(map[string]*entities.StreamEvent)}
	if len(changes.EventIds) == 0 {
		return res, nil
	}

	// PostgreSQL doesn't explicitly limit the number of arguments,
	// but some drivers may set a limit of 32767 bind arguments.
	// Because each Primary Key needs 2 bind arguments for each condition,
	// we are totally safe to use all Primary Keys at once. No need to chunk.
	var names = make([]string, len(changes.EventIds))
	var args = pgx.NamedArgs{}
	for i, id := range changes.EventIds {
		binding := fmt.Sprintf("event_id_%d", i)
		names[i] = fmt.Sprintf("@%s", binding)
		args[binding] = id
	}

	table := pgx.Identifier{entities.CollectionStreamEvent(req.Consumer.StreamName)}.Sanitize()
	query := fmt.Sprintf(ConsumerJobPullSQL, table, strings.Join(names, ","))
	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event entities.StreamEvent
		err := rows.Scan(
			&event.Topic,
			&event.EventId,
			&event.Body,
			&event.Metadata,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		res.Events[event.EventId] = &event
	}

	return res, nil
}
