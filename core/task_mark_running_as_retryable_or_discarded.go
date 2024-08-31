package core

import (
	"context"
	_ "embed"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/xvalidator"
)

//go:embed task_mark_running_as_retryable_or_discarded.sql
var TaskMarkRunningAsRetryableOrDiscardedSql string

type TaskMarkRunningAsRetryableOrDiscardedReq struct {
	Consumer *entities.ConsumerRegistry `validate:"required"`
	Tasks    []*entities.Task           `validate:"required,gt=0,dive,required"`
	Error    entities.AttemptedError    `validate:"required"`
}

type TaskMarkRunningAsRetryableOrDiscardedRes struct {
	Updated []string
	Noop    []string
	States  map[string]entities.TaskState
}

func (req *TaskMarkRunningAsRetryableOrDiscardedReq) Do(ctx context.Context, tx pgx.Tx) (*TaskMarkRunningAsRetryableOrDiscardedRes, error) {
	err := xvalidator.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	modified := make(map[string]bool)

	var names = make([]string, len(req.Tasks))
	var args = pgx.NamedArgs{
		"attempt_max":     req.Consumer.AttemptMax,
		"finalized_at":    time.Now().UnixMilli(),
		"attempted_error": req.Error,
		"discarded_state": int(entities.StateDiscarded),
		"retryable_state": int(entities.StateRetryable),
		"running_state":   int(entities.StateRunning),
	}
	for i, task := range req.Tasks {
		// we assume that tasks are not able to update firstly
		modified[task.EventId] = false

		bind := fmt.Sprintf("event_id_%d", i)
		names[i] = fmt.Sprintf("@%s", fmt.Sprintf("event_id_%d", i))
		args[bind] = task.EventId
	}

	table := pgx.Identifier{entities.Collection(req.Consumer.Id)}.Sanitize()
	query := fmt.Sprintf(TaskMarkRunningAsRetryableOrDiscardedSql, table, strings.Join(names, ","))

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	states := make(map[string]entities.TaskState)
	for rows.Next() {
		var eventId string
		var state int16
		if err := rows.Scan(&eventId, &state); err != nil {
			return nil, err
		}

		// report that we was able to update the task
		modified[eventId] = true
		states[eventId] = entities.TaskState(state)
	}

	// rows.Err returns any error that occurred while reading
	// always check it before finishing the read
	if err := rows.Err(); err != nil {
		return nil, err
	}

	res := &TaskMarkRunningAsRetryableOrDiscardedRes{
		Updated: make([]string, 0),
		Noop:    make([]string, 0),
		States:  states,
	}
	for id, updated := range modified {
		if updated {
			res.Updated = append(res.Updated, id)
			continue
		}

		res.Noop = append(res.Noop, id)
	}

	return res, nil
}
