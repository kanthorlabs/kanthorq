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

//go:embed task_convert.sql
var TaskConvertSql string

type TaskConvertReq struct {
	Consumer     *entities.ConsumerRegistry `validate:"required"`
	EventIds     []string                   `validate:"required,gt=0,lte=500,dive,required"`
	InitialState entities.TaskState         `validate:"is_enum"`
}

type TaskConvertRes struct {
	Tasks    map[string]*entities.Task
	EventIds []string
}

func (req *TaskConvertReq) Do(ctx context.Context, tx pgx.Tx) (*TaskConvertRes, error) {
	err := xvalidator.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	var args = pgx.NamedArgs{
		"intial_state": int16(req.InitialState),
		"schedule_at":  time.Now().UnixMilli() + req.Consumer.VisibilityTimeout,
	}
	var names = make([]string, len(req.EventIds))
	for i, id := range req.EventIds {
		bind := fmt.Sprintf("event_id_%d", i)
		names[i] = fmt.Sprintf("@%s", bind)
		args[bind] = id
	}

	stable := pgx.Identifier{entities.Collection(req.Consumer.StreamId)}.Sanitize()
	ctable := pgx.Identifier{entities.Collection(req.Consumer.Id)}.Sanitize()
	query := fmt.Sprintf(TaskConvertSql, ctable, stable, strings.Join(names, ","))

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := &TaskConvertRes{Tasks: make(map[string]*entities.Task)}
	for rows.Next() {
		var task entities.Task
		err := rows.Scan(
			&task.EventId,
			&task.Subject,
			&task.State,
			&task.ScheduleAt,
			&task.FinalizedAt,
			&task.AttemptCount,
			&task.AttemptedAt,
			&task.AttemptedError,
			&task.Metadata,
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
