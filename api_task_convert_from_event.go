package kanthorq

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/pkg/validator"
)

//go:embed api_task_convert_from_event.sql
var TaskConvertFromEventSql string

type TaskConvertFromEventReq struct {
	ConsumerName     string    `validate:"required,is_collection_name"`
	Size             int       `validate:"required,gt=0"`
	InitialTaskState TaskState `validate:"required,is_enum"`
}

type TaskConvertFromEventRes struct {
	Tasks []*Task
}

func (req *TaskConvertFromEventReq) Do(ctx context.Context, tx pgx.Tx) (*TaskConvertFromEventRes, error) {
	err := validator.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	// lock consumer firstly
	consumer, err := req.lock(ctx, tx)
	if err != nil {
		return nil, err
	}

	// convert event into task
	args := pgx.NamedArgs{
		"consumer_topic":  consumer.Topic,
		"consumer_cursor": consumer.Cursor,
		"size":            req.Size,
		"intial_state":    int(req.InitialTaskState),
	}
	ctable := pgx.Identifier{Collection(consumer.Name)}.Sanitize()
	stable := pgx.Identifier{Collection(consumer.StreamName)}.Sanitize()
	query := fmt.Sprintf(TaskConvertFromEventSql, ctable, stable)

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := &TaskConvertFromEventRes{Tasks: make([]*Task, 0, req.Size)}
	for rows.Next() {
		var task Task
		err = rows.Scan(
			&task.EventId,
			&task.Topic,
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
		res.Tasks = append(res.Tasks, &task)
	}

	return res, nil
}

func (req *TaskConvertFromEventReq) lock(ctx context.Context, tx pgx.Tx) (*ConsumerRegistry, error) {
	var args = pgx.NamedArgs{
		"consumer_name": req.ConsumerName,
	}
	var consumer ConsumerRegistry
	var err = tx.QueryRow(ctx, ConsumerLockSql, args).Scan(
		&consumer.Name,
		&consumer.StreamName,
		&consumer.Topic,
		&consumer.Cursor,
		&consumer.AttemptMax,
		&consumer.CreatedAt,
		&consumer.UpdatedAt,
	)

	return &consumer, err
}
