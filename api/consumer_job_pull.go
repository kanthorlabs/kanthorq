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

// NewConsumerJobPull is the main function of this system
// it will find available jobs, make it become running, and return the events themself
func NewConsumerJobPull(consumer *entities.Consumer, size int, vt time.Duration) *ConsumerJobPullReq {
	return &ConsumerJobPullReq{
		Consumer:          consumer,
		Size:              size,
		VisibilityTimeout: vt,
		FromState:         entities.StateAvailable,
		ToState:           entities.StateRunning,
	}
}

// NewConsumerJobRetry will find retryable jobs, make them become running, and return the events themself
func NewConsumerJobRetry(consumer *entities.Consumer, size int, vt time.Duration) *ConsumerJobPullReq {
	return &ConsumerJobPullReq{
		Consumer:          consumer,
		Size:              size,
		VisibilityTimeout: vt,
		FromState:         entities.StateRetryable,
		ToState:           entities.StateRunning,
	}
}

//go:embed consumer_job_pull.sql
var ConsumerJobPullSQL string

type ConsumerJobPullReq struct {
	Consumer          *entities.Consumer
	Size              int
	VisibilityTimeout time.Duration
	FromState         entities.JobState
	ToState           entities.JobState
}

type ConsumerJobPullRes struct {
	Events []*entities.StreamEvent
}

func (req *ConsumerJobPullReq) Do(ctx context.Context, tx pgx.Tx) (*ConsumerJobPullRes, error) {
	command := NewConsumerJobStateChange(
		req.Consumer,
		req.Size,
		req.VisibilityTimeout,
		req.FromState,
		req.ToState,
	)
	changes, err := command.Do(ctx, tx)
	if err != nil {
		return nil, err
	}

	res := &ConsumerJobPullRes{Events: make([]*entities.StreamEvent, 0)}
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
		res.Events = append(res.Events, &event)
	}

	return res, nil
}
