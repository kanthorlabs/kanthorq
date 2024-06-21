package api

import (
	"context"
	"errors"
	"fmt"
	"time"

	_ "embed"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
)

func NewConsumerJobStateChange(
	consumer *entities.Consumer,
	size int,
	vt time.Duration,
	fromState entities.JobState,
	toState entities.JobState,
) *ConsumerJobStateChangeReq {
	return &ConsumerJobStateChangeReq{
		Consumer:          consumer,
		Size:              size,
		VisibilityTimeout: vt,
		FromState:         fromState,
		ToState:           toState,
	}
}

//go:embed consumer_job_state_change.sql
var ConsumerJobStateChangeSQL string

type ConsumerJobStateChangeReq struct {
	Consumer          *entities.Consumer
	Size              int
	VisibilityTimeout time.Duration
	FromState         entities.JobState
	ToState           entities.JobState
}

type ConsumerJobStateChangeRes struct {
	EventIds []string
}

func (req *ConsumerJobStateChangeReq) Do(ctx context.Context, tx pgx.Tx) (*ConsumerJobStateChangeRes, error) {
	res := &ConsumerJobStateChangeRes{}

	args := pgx.NamedArgs{
		"attempt_at":       time.Now().UnixMilli(),
		"from_state":       int(req.FromState),
		"to_state":         int(req.ToState),
		"size":             req.Size,
		"next_schedule_at": time.Now().Add(req.VisibilityTimeout).UnixMilli(),
	}

	table := pgx.Identifier{entities.CollectionConsumerJob(req.Consumer.Name)}.Sanitize()
	query := fmt.Sprintf(ConsumerJobStateChangeSQL, table, table)
	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		res.EventIds = append(res.EventIds, id)
	}

	// if error is nil or error is pgx.ErrNoRows, we accept it
	if err := rows.Err(); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	return res, nil
}
