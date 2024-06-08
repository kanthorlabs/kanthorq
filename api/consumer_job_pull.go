package api

import (
	"context"
	_ "embed"
	"errors"
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
	if len(changes.PrimaryKeys) == 0 {
		return res, err
	}

	// PostgreSQL doesn't explicitly limit the number of arguments,
	// but some drivers may set a limit of 32767 bind arguments.
	// Because each Primary Key needs 2 bind arguments for each condition,
	// we are totally safe to use all Primary Keys at once. No need to chunk.
	var names = make([]string, len(changes.PrimaryKeys))
	var args = pgx.NamedArgs{}
	for i, pk := range changes.PrimaryKeys {
		topicBind := fmt.Sprintf("topic_%d", i)
		eventIdBind := fmt.Sprintf("event_id_%d", i)

		names[i] = fmt.Sprintf("(@%s, @%s)", topicBind, eventIdBind)
		args[topicBind] = pk.Topic
		args[eventIdBind] = pk.EventId
	}

	table := pgx.Identifier{entities.CollectionStreamEvent(req.Consumer.StreamName)}.Sanitize()
	query := fmt.Sprintf(ConsumerJobPullSQL, table, strings.Join(names, ","))
	rows, err := tx.Query(ctx, query, args)
	// in the initial state, there is no rows to pull
	// so pgx will return ErrNoRows, need to cast it as successful case
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return res, nil
	}
	defer rows.Close()

	for rows.Next() {
		var event entities.StreamEvent
		if err := rows.Scan(&event.Topic, &event.EventId, &event.CreatedAt); err != nil {
			return nil, err
		}
		res.Events[event.EventId] = &event
	}

	return res, nil
}
