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

//go:embed task_mark_cancelled.sql
var TaskMarkCancelledSql string

type TaskMarkCancelledReq struct {
	Consumer *entities.ConsumerRegistry `validate:"required"`
	Tasks    []*entities.Task           `validate:"required,gt=0,lte=500,dive,required"`
}

type TaskMarkCancelledRes struct {
	Updated []string
	Noop    []string
}

func (req *TaskMarkCancelledReq) Do(ctx context.Context, tx pgx.Tx) (*TaskMarkCancelledRes, error) {
	err := xvalidator.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	modified := make(map[string]bool)

	// As you already know IN expects a list of scalar expressions, not an array,
	// however pgtype.Int4Array represents an array, not a list of scalar expressions.
	// I tested and confirmed this is the case.
	// Using ANY($1) in the WHERE clause is a workaround but it will use filter instead of index scan
	// that why we use IN instead of ANY
	var names = make([]string, len(req.Tasks))
	var args = pgx.NamedArgs{
		"cancelled_state": int16(entities.StateCancelled),
		"pending_state":   int16(entities.StatePending),
		"available_state": int16(entities.StateAvailable),
		"retryable_state": int16(entities.StateRetryable),
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
	query := fmt.Sprintf(TaskMarkCancelledSql, table, strings.Join(names, ","))
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

	res := &TaskMarkCancelledRes{Updated: make([]string, 0), Noop: make([]string, 0)}
	for id, updated := range modified {
		if updated {
			res.Updated = append(res.Updated, id)
			continue
		}

		res.Noop = append(res.Noop, id)
	}

	return res, nil
}
