package kanthorq

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func StreamPutEvents(ctx context.Context, req *StreamPutEventsReq, conn *pgx.Conn) (*StreamPutEventsRes, error) {
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	res, err := req.Do(ctx, tx)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return res, nil
}

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
			event.Topic,
			event.Body,
			event.Metadata,
			event.CreatedAt,
		}
	}

	identifier := pgx.Identifier{Collection(req.Stream.Name)}
	inserted, err := tx.CopyFrom(ctx, identifier, EventProps(), pgx.CopyFromRows(rows))
	if err != nil {
		return nil, err
	}

	return &StreamPutEventsRes{InsertCount: inserted}, nil
}
