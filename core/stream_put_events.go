package core

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/xvalidator"
)

type StreamPutEventsReq struct {
	Stream *entities.StreamRegistry `validate:"required"`
	Events []*entities.Event        `validate:"required"`
}

type StreamPutEventsRes struct {
	InsertCount int64
}

func (req *StreamPutEventsReq) Do(ctx context.Context, tx pgx.Tx) (*StreamPutEventsRes, error) {
	err := xvalidator.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

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

	identifier := pgx.Identifier{entities.Collection(req.Stream.Id)}
	inserted, err := tx.CopyFrom(ctx, identifier, entities.Properties(entities.Event{}), pgx.CopyFromRows(rows))
	if err != nil {
		return nil, err
	}

	return &StreamPutEventsRes{InsertCount: inserted}, nil
}
