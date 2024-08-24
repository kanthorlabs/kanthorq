package core

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/xvalidator"
)

//go:embed task_state_transition.sql
var TaskStateTransitionSql string

type TaskStateTransitionReq struct {
	Consumer  *entities.ConsumerRegistry `validate:"required"`
	FromState entities.TaskState         `validate:"required,is_enum"`
	ToState   entities.TaskState         `validate:"required,is_enum"`
	Size      int                        `validate:"required,gt=0"`
}

type TaskStateTransitionRes struct {
	Tasks    map[string]*entities.Task
	EventIds []string
}

func (req *TaskStateTransitionReq) Do(ctx context.Context, tx pgx.Tx) (*TaskStateTransitionRes, error) {
	err := xvalidator.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	now := time.Now().UnixMilli()
	var args = pgx.NamedArgs{
		"attempt_max":  req.Consumer.AttemptMax,
		"attempted_at": now,
		"schedule_at":  now + req.Consumer.VisibilityTimeout,
		"from_state":   int(req.FromState),
		"to_state":     int(req.ToState),
		"size":         req.Size,
	}
	table := pgx.Identifier{entities.Collection(req.Consumer.Id)}.Sanitize()
	query := fmt.Sprintf(TaskStateTransitionSql, table, table)

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := &TaskStateTransitionRes{
		Tasks: make(map[string]*entities.Task),
	}
	for rows.Next() {
		var task entities.Task
		err := rows.Scan(
			&task.EventId,
			&task.Subject,
			&task.State,
			&task.ScheduleAt,
			&task.AttemptCount,
			&task.AttemptedAt,
			&task.FinalizedAt,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		res.Tasks[task.EventId] = &task
		res.EventIds = append(res.EventIds, task.EventId)
	}

	// rows.Err returns any error that occurred while reading
	// always check it before finishing the read
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
