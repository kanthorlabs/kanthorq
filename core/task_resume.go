package core

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/xvalidator"
)

//go:embed task_resume.sql
var TaskResumeSql string

type TaskResumeReq struct {
	Consumer *entities.ConsumerRegistry `validate:"required"`
	Tasks    []*entities.Task           `validate:"required,gt=0,lte=500,dive,required"`
}

type TaskResumeRes struct {
	Updated []string
	Noop    []string
}

func (req *TaskResumeReq) Do(ctx context.Context, tx pgx.Tx) (*TaskResumeRes, error) {
	err := xvalidator.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	modified := make(map[string]bool)

	var names = make([]string, len(req.Tasks))
	var args = pgx.NamedArgs{
		"state_running":   int16(entities.StateRunning),
		"state_discarded": int16(entities.StateDiscarded),
		"state_cancelled": int16(entities.StateCancelled),
	}

	for i, task := range req.Tasks {
		// we assume that tasks are not able to update firstly
		modified[task.EventId] = false

		bind := fmt.Sprintf("event_id_%d", i)
		names[i] = fmt.Sprintf("@%s", bind)
		args[bind] = task.EventId
	}

	table := pgx.Identifier{entities.Collection(req.Consumer.Id)}.Sanitize()
	query := fmt.Sprintf(TaskResumeSql, table, strings.Join(names, ","))
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

	res := &TaskResumeRes{Updated: make([]string, 0), Noop: make([]string, 0)}
	for id, updated := range modified {
		if updated {
			res.Updated = append(res.Updated, id)
			continue
		}

		res.Noop = append(res.Noop, id)
	}

	return res, nil
}
