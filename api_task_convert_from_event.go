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
	Consumer         *ConsumerRegistry `validate:"required"`
	Size             int               `validate:"required,gt=0"`
	InitialTaskState TaskState         `validate:"required,is_enum"`
}

type TaskConvertFromEventRes struct {
	EventIds   []string
	NextCursor string
	Tasks      map[string]*Task
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
	ctable := pgx.Identifier{Collection(consumer.Id)}.Sanitize()
	stable := pgx.Identifier{Collection(consumer.StreamId)}.Sanitize()
	query := fmt.Sprintf(TaskConvertFromEventSql, ctable, stable)

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := &TaskConvertFromEventRes{Tasks: make(map[string]*Task)}
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

		res.EventIds = append(res.EventIds, task.EventId)
		// always overwrite latst event id as cursor
		res.NextCursor = task.EventId
		res.Tasks[task.EventId] = &task
	}

	// rows.Err returns any error that occurred while reading
	// always check it before finishing the read
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// update cursor
	if err := req.update(ctx, tx, res.NextCursor); err != nil {
		return nil, err
	}

	return res, nil
}

func (req *TaskConvertFromEventReq) lock(ctx context.Context, tx pgx.Tx) (*ConsumerRegistry, error) {
	var args = pgx.NamedArgs{
		"consumer_name": req.Consumer.Name,
	}
	var consumer ConsumerRegistry
	var err = tx.QueryRow(ctx, ConsumerLockSql, args).Scan(
		&consumer.StreamId,
		&consumer.StreamName,
		&consumer.Id,
		&consumer.Name,
		&consumer.Topic,
		&consumer.Cursor,
		&consumer.AttemptMax,
		&consumer.CreatedAt,
		&consumer.UpdatedAt,
	)

	return &consumer, err
}

func (req *TaskConvertFromEventReq) update(ctx context.Context, tx pgx.Tx, cursor string) error {
	var args = pgx.NamedArgs{
		"consumer_name":   req.Consumer.Name,
		"consumer_cursor": cursor,
	}

	_, err := tx.Exec(ctx, ConsumerUpdateCursorSql, args)
	return err
}
