package core

import (
	"context"
	_ "embed"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/validator"
)

//go:embed task_mark_running_as_completed.sql
var TaskMarkRunningAsCompletedSql string

type TaskMarkRunningAsCompletedReq struct {
	Consumer *entities.ConsumerRegistry `validate:"required"`
	Tasks    []*entities.Task           `validate:"required,gt=0,dive,required"`
}

type TaskMarkRunningAsCompletedRes struct {
	Updated []string
	Noop    []string
}

func (req *TaskMarkRunningAsCompletedReq) Do(ctx context.Context, tx pgx.Tx) (*TaskMarkRunningAsCompletedRes, error) {
	err := validator.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	modified := make(map[string]bool)

	var names = make([]string, len(req.Tasks))
	var args = pgx.NamedArgs{
		"running_state":   int(entities.StateRunning),
		"completed_state": int(entities.StateCompleted),
		"finalized_at":    time.Now().UnixMilli(),
	}
	for i, task := range req.Tasks {
		// we assume that tasks are not able to update firstly
		modified[task.EventId] = false

		bind := fmt.Sprintf("event_id_%d", i)
		names[i] = fmt.Sprintf("@%s", bind)
		args[bind] = task.EventId
	}

	table := pgx.Identifier{entities.Collection(req.Consumer.Id)}.Sanitize()
	query := fmt.Sprintf(TaskMarkRunningAsCompletedSql, table, strings.Join(names, ","))
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

	res := &TaskMarkRunningAsCompletedRes{Updated: make([]string, 0), Noop: make([]string, 0)}
	for id, updated := range modified {
		if updated {
			res.Updated = append(res.Updated, id)
			continue
		}

		res.Noop = append(res.Noop, id)
	}

	return res, nil
}
