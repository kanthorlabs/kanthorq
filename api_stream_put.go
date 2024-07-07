package kanthorq

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func StreamPut(stream *StreamRegistry, events []*Event) *StreamPutReq {
	return &StreamPutReq{Stream: stream, Events: events}
}

type StreamPutReq struct {
	Stream *StreamRegistry
	Events []*Event
}

type StreamPutRes struct {
	InsertCount int64
}

func (req *StreamPutReq) Do(ctx context.Context, tx pgx.Tx) (*StreamPutRes, error) {
	if len(req.Events) == 0 {
		return &StreamPutRes{InsertCount: 0}, nil
	}

	var rows = make([][]any, len(req.Events))
	for i, event := range req.Events {
		rows[i] = []any{
			event.Id,
			event.Topic,
			event.Body,
			event.Metadata,
			event.CreatedAt,
		}
	}

	identifier := pgx.Identifier{StreamCollection(req.Stream.Name)}
	inserted, err := tx.CopyFrom(ctx, identifier, EventProps, pgx.CopyFromRows(rows))
	if err != nil {
		return nil, err
	}

	return &StreamPutRes{InsertCount: inserted}, nil
}
