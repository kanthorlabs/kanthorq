package core

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/xvalidator"
)

type ConsumerPutTasksReq struct {
	Consumer *entities.ConsumerRegistry `validate:"required"`
	Tasks    []*entities.Task           `validate:"required,dive,required"`
}

type ConsumerPutTasksRes struct {
	InsertCount int64
}

func (req *ConsumerPutTasksReq) Do(ctx context.Context, tx pgx.Tx) (*ConsumerPutTasksRes, error) {
	err := xvalidator.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	var rows = make([][]any, len(req.Tasks))
	for i, task := range req.Tasks {
		rows[i] = []any{
			task.EventId,
			task.Subject,
			task.State,
			task.ScheduleAt,
			task.AttemptCount,
			task.AttemptedAt,
			task.AttemptedError,
			task.FinalizedAt,
			task.Metadata,
			task.CreatedAt,
			task.UpdatedAt,
		}
	}

	identifier := pgx.Identifier{entities.Collection(req.Consumer.Id)}
	inserted, err := tx.CopyFrom(ctx, identifier, entities.Properties(entities.Task{}), pgx.CopyFromRows(rows))
	if err != nil {
		return nil, err
	}

	return &ConsumerPutTasksRes{InsertCount: inserted}, nil
}
