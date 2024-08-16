package kanthorq

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type ConsumerPutTasksReq struct {
	Consumer *ConsumerRegistry
	Tasks    []*Task
}

type ConsumerPutTasksRes struct {
	InsertCount int64
}

func (req *ConsumerPutTasksReq) Do(ctx context.Context, tx pgx.Tx) (*ConsumerPutTasksRes, error) {
	if len(req.Tasks) == 0 {
		return &ConsumerPutTasksRes{InsertCount: 0}, nil
	}

	var rows = make([][]any, len(req.Tasks))
	for i, event := range req.Tasks {
		rows[i] = []any{
			event.EventId,
			event.Subject,
			event.State,
			event.ScheduleAt,
			event.AttemptCount,
			event.AttemptedAt,
			event.FinalizedAt,
			event.CreatedAt,
			event.UpdatedAt,
		}
	}

	identifier := pgx.Identifier{Collection(req.Consumer.Id)}
	inserted, err := tx.CopyFrom(ctx, identifier, Properties(Task{}), pgx.CopyFromRows(rows))
	if err != nil {
		return nil, err
	}

	return &ConsumerPutTasksRes{InsertCount: inserted}, nil
}
