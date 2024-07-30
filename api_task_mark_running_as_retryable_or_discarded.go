package kanthorq

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/pkg/validator"
)

//go:embed api_task_mark_running_as_retryable_or_discarded.sql
var TaskMarkRunningAsRetryableOrDiscardedSql string

type TaskMarkRunningAsRetryableOrDiscardedReq struct {
	Consumer *ConsumerRegistry `validate:"required"`
	Tasks    []*Task           `validate:"required,gt=0,dive,required"`
}

type TaskMarkRunningAsRetryableOrDiscardedRes struct {
	Updated []string
	Noop    []string
}

func (req *TaskMarkRunningAsRetryableOrDiscardedReq) Do(ctx context.Context, tx pgx.Tx) (*TaskMarkRunningAsRetryableOrDiscardedRes, error) {
	err := validator.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	modified := make(map[string]bool)

	var names = make([]string, len(req.Tasks))
	var args = pgx.NamedArgs{
		"attempt_max":     req.Consumer.AttemptMax,
		"discarded_state": int(StateDiscarded),
		"retryable_state": int(StateRetryable),
		"running_state":   int(StateRunning),
	}
	for i, task := range req.Tasks {
		// we assume that tasks are not able to update firstly
		modified[task.EventId] = false

		bind := fmt.Sprintf("event_id_%d", i)
		names[i] = fmt.Sprintf("@%s", fmt.Sprintf("event_id_%d", i))
		args[bind] = task.EventId
	}

	table := pgx.Identifier{Collection(req.Consumer.Name)}.Sanitize()
	query := fmt.Sprintf(TaskMarkRunningAsRetryableOrDiscardedSql, table, strings.Join(names, ","))

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var eventId string
		if err := rows.Scan(&eventId); err != nil {
			return nil, err
		}

		// report that we was able to update the task
		modified[eventId] = true
	}

	// rows.Err returns any error that occurred while reading
	// always check it before finishing the read
	if err := rows.Err(); err != nil {
		return nil, err
	}

	res := &TaskMarkRunningAsRetryableOrDiscardedRes{Updated: make([]string, 0), Noop: make([]string, 0)}
	for id, updated := range modified {
		if updated {
			res.Updated = append(res.Updated, id)
			continue
		}

		res.Noop = append(res.Noop, id)
	}

	return res, nil
}
