package api

import (
	"context"
	"fmt"
	"time"

	_ "embed"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
)

func NewConsumerJobStateChange(
	consumer *entities.Consumer,
	size int,
	fromState entities.JobState,
	toState entities.JobState,
	vt time.Duration,
) *ConsumerJobStateChangeReq {
	return &ConsumerJobStateChangeReq{
		Consumer:          consumer,
		Size:              size,
		FromState:         fromState,
		ToState:           toState,
		VisibilityTimeout: vt,
	}
}

//go:embed consumer_job_state_change.sql
var ConsumerJobStateChangeSQL string

type ConsumerJobStateChangeReq struct {
	Consumer          *entities.Consumer
	Size              int
	FromState         entities.JobState
	ToState           entities.JobState
	VisibilityTimeout time.Duration
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

	return res, rows.Err()
}
