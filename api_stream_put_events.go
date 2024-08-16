package kanthorq

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type StreamPutEventsReq struct {
	Stream *StreamRegistry
	Events []*Event
}

type StreamPutEventsRes struct {
	InsertCount int64
}

func (req *StreamPutEventsReq) Do(ctx context.Context, tx pgx.Tx) (*StreamPutEventsRes, error) {
	if len(req.Events) == 0 {
		return &StreamPutEventsRes{InsertCount: 0}, nil
	}

	var rows = make([][]any, len(req.Events))
	for i, event := range req.Events {
		rows[i] = []any{
			event.Id,
			event.Subject,
			event.Body,
			event.Metadata,
			event.CreatedAt,
		}
	}

	identifier := pgx.Identifier{Collection(req.Stream.Id)}
	inserted, err := tx.CopyFrom(ctx, identifier, Properties(Event{}), pgx.CopyFromRows(rows))
	if err != nil {
		return nil, err
	}

	return &StreamPutEventsRes{InsertCount: inserted}, nil
}
