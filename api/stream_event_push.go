package api

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
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
	ctx, span := telemetry.Tracer.Start(ctx, "api.StreamEventPush", trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	if len(req.Events) == 0 {
		span.SetAttributes(attribute.Bool("ERROR.STREAM_EVENTS.EMPTY", true))
		return nil, fmt.Errorf("ERROR.STREAM_EVENTS.EMPTY: %s", req.Stream.Name)
	}

	var entries = make([][]any, len(req.Events))
	for i, event := range req.Events {
		// inject tracing context into metadata
		carrier := propagation.MapCarrier{}
		otel.GetTextMapPropagator().Inject(ctx, carrier)
		for k, v := range carrier {
			event.Metadata[k] = v
		}

		entries[i] = []any{
			event.EventId,
			event.Topic,
			event.Body,
			event.Metadata,
			event.CreatedAt,
		}
	}
	span.SetAttributes(attribute.Int("event_count", len(entries)))

	inserted, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{entities.CollectionStreamEvent(req.Stream.Name)},
		entities.CollectionStreamEventProps,
		pgx.CopyFromRows(entries),
	)
	if err != nil {
		span.RecordError(err, trace.WithStackTrace(true))
		return nil, err
	}

	return &StreamEventPushRes{inserted}, nil
}
