package api

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
)

func StreamEventPush(stream *entities.Stream, events []*entities.StreamEvent) *StreamEventPushReq {
	return &StreamEventPushReq{
		Stream: stream,
		Events: events,
	}
}

type StreamEventPushReq struct {
	Stream *entities.Stream
	Events []*entities.StreamEvent
}

type StreamEventPushRes struct {
	InsertCount int64
}

func (req *StreamEventPushReq) Do(ctx context.Context, tx pgx.Tx) (*StreamEventPushRes, error) {
	if len(req.Events) == 0 {
		return nil, fmt.Errorf("ERROR.STREAM_EVENTS.EMPTY: %s", req.Stream.Name)
	}

	var entries = make([][]any, len(req.Events))
	for i, event := range req.Events {
		entries[i] = []any{event.Topic, event.EventId, event.CreatedAt}
	}

	inserted, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{entities.CollectionStreamEvent(req.Stream.Name)},
		entities.CollectionStreamEventProps,
		pgx.CopyFromRows(entries),
	)
	if err != nil {
		return nil, err
	}

	return &StreamEventPushRes{inserted}, nil
}
