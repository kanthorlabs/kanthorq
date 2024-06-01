package api

import (
	"context"
	_ "embed"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/common/clock"
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

func (req *ConsumerJobPullReq) Do(ctx context.Context, tx pgx.Tx, clock clock.Clock) (*ConsumerJobPullRes, error) {
	changes, err := req.changes().Do(ctx, tx, clock)
	if err != nil {
		return nil, err
	}

	res := &ConsumerJobPullRes{Events: make(map[string]*entities.StreamEvent)}
	if len(changes.PrimaryKeys) == 0 {
		return res, err
	}

	// PostgreSQL doesn't explicitly limit the number of arguments,
	// but some drivers may set a limit of 32767 bind arguments
	// Because our Primary Key need 2 bind arguments for each condition
	// we are totally safe to use all Primary Keys at once. no need to chunking
	var names = make([]string, len(changes.PrimaryKeys))
	var args = pgx.NamedArgs{}
	for i, pk := range changes.PrimaryKeys {
		topic := fmt.Sprintf("topic_%d", i)
		eventId := fmt.Sprintf("event_id_%d", i)

		names[i] = fmt.Sprintf("(@%s, @%s)", topic, eventId)
		args[topic] = pk.Topic
		args[eventId] = pk.EventId
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
		if err := rows.Scan(&event.Topic, &event.EventId, &event.CreatedAt); err != nil {
			return nil, err
		}
		res.Events[event.EventId] = &event
	}

	return res, nil
}

func (req *ConsumerJobPullReq) changes() *ConsumerJobStateChangeReq {
	return &ConsumerJobStateChangeReq{
		Consumer:          req.Consumer,
		Size:              req.Size,
		FromState:         entities.StateAvailable,
		ToState:           entities.StateRunning,
		VisibilityTimeout: req.VisibilityTimeout,
	}
}
