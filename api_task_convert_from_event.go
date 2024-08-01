package kanthorq

import (
	"context"
	_ "embed"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/pkg/idx"
	"github.com/kanthorlabs/kanthorq/pkg/utils"
	"github.com/kanthorlabs/kanthorq/pkg/validator"
)

//go:embed api_task_convert_from_event.sql
var TaskConvertFromEventSql string

type TaskConvertFromEventReq struct {
	Consumer         *ConsumerRegistry `validate:"required"`
	InitialTaskState TaskState         `validate:"required,is_enum"`
	Size             int               `validate:"required,gt=0"`
	ScanWindow       int64             `validate:"required,gte=1000"`
	ScanRoundMax     int               `validate:"required,gt=0"`
	ScanRoundDelay   int64             `validate:"required,gte=1000"`
}

type TaskConvertFromEventRes struct {
	EventIds []string
	Tasks    map[string]*Task
}

func (req *TaskConvertFromEventReq) Do(ctx context.Context, tx pgx.Tx) (*TaskConvertFromEventRes, error) {
	err := validator.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	// lock consumer firstly
	if err := req.lock(ctx, tx); err != nil {
		return nil, err
	}

	var res = &TaskConvertFromEventRes{Tasks: make(map[string]*Task)}
	var round int

	for len(res.EventIds) < req.Size && round < req.ScanRoundMax {
		round++

		out, err := req.scan(ctx, tx)
		if err != nil {
			// @TODO: log the error here
			break
		}

		if len(out.Tasks) > 0 {
			res.EventIds = append(res.EventIds, out.EventIds...)
			utils.MergeMaps(res.Tasks, out.Tasks)
			continue
		}

		// @TODO: log no task was found here
		time.Sleep(time.Millisecond * time.Duration(req.ScanRoundDelay))
	}

	// updating
	if err := req.update(ctx, tx); err != nil {
		return nil, err
	}

	return res, nil
}

func (req *TaskConvertFromEventReq) lock(ctx context.Context, tx pgx.Tx) error {
	var args = pgx.NamedArgs{
		"consumer_name": req.Consumer.Name,
	}

	var consumer ConsumerRegistry
	err := tx.QueryRow(ctx, ConsumerLockSql, args).Scan(
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
	if err != nil {
		return err
	}

	// override consumer with the locked consumer
	req.Consumer = &consumer
	return nil
}

func (req *TaskConvertFromEventReq) scan(ctx context.Context, tx pgx.Tx) (*TaskConvertFromEventRes, error) {
	res := &TaskConvertFromEventRes{Tasks: make(map[string]*Task)}

	args := pgx.NamedArgs{
		"intial_state":          int(req.InitialTaskState),
		"consumer_topic_filter": TopicFilter(req.Consumer.Topic),
		"size":                  req.Size,
	}

	where := []string{
		"id > @consumer_cursor_start",
	}

	// match all topics so we don't need to filter by topic
	if req.Consumer.Topic == TopicAll {
		args["consumer_cursor_start"] = req.Consumer.Cursor
	} else {
		where = append(where, "topic LIKE @consumer_topic_filter")
		args["consumer_topic_filter"] = TopicFilter(req.Consumer.Topic)

		// starting with fresh consumer, cursor will be empty
		if req.Consumer.Cursor == "" {
			args["consumer_cursor_start"] = ""
		} else {
			where = append(where, "id < @consumer_cursor_end")
			args["consumer_cursor_end"] = idx.Next(req.Consumer.Cursor, time.Millisecond*time.Duration(req.ScanWindow))
		}
	}

	ctable := pgx.Identifier{Collection(req.Consumer.Id)}.Sanitize()
	stable := pgx.Identifier{Collection(req.Consumer.StreamId)}.Sanitize()
	query := fmt.Sprintf(TaskConvertFromEventSql, ctable, stable, strings.Join(where, " AND "))

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
		res.Tasks[task.EventId] = &task

		// override locked consumer cursor with latest event id
		req.Consumer.Cursor = task.EventId
	}

	// rows.Err returns any error that occurred while reading
	// always check it before finishing the read
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (req *TaskConvertFromEventReq) update(ctx context.Context, tx pgx.Tx) error {
	var args = pgx.NamedArgs{
		"consumer_name":   req.Consumer.Name,
		"consumer_cursor": req.Consumer.Cursor,
	}
	_, err := tx.Exec(ctx, ConsumerUpdateCursorSql, args)
	return err
}
