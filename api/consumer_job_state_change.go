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

func ConsumerJobStateChange(
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
	PrimaryKeys []*entities.EventPrimaryKey
}

func (req *ConsumerJobStateChangeReq) Do(ctx context.Context, tx pgx.Tx) (*ConsumerJobStateChangeRes, error) {
	res := &ConsumerJobStateChangeRes{}

	args := pgx.NamedArgs{
		"attempt_at":       time.Now().UnixMilli(),
		"from_state":       req.FromState,
		"to_state":         req.ToState,
		"size":             req.Size,
		"next_schedule_at": time.Now().Add(req.VisibilityTimeout).UnixMilli(),
	}

	table := pgx.Identifier{entities.CollectionConsumerJob(req.Consumer.Name)}.Sanitize()
	query := fmt.Sprintf(ConsumerJobStateChangeSQL, table, table)
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
		var pk entities.EventPrimaryKey
		err = rows.Scan(&pk.Topic, &pk.EventId)
		if err != nil {
			return nil, err
		}
		res.PrimaryKeys = append(res.PrimaryKeys, &pk)
	}

	return res, rows.Err()
}
