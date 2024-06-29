package q

import (
	"context"
	_ "embed"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func NewConsumerJobMarkCancelled(consumer *entities.Consumer, events []*entities.StreamEvent) *ConsumerJobMarkCancelledReq {
	return &ConsumerJobMarkCancelledReq{
		Consumer: consumer,
		Events:   events,
	}
}

//go:embed consumer_job_mark_cancelled.sql
var ConsumerJobMarkCancelledSQL string

type ConsumerJobMarkCancelledReq struct {
	Consumer *entities.Consumer
	Events   []*entities.StreamEvent
}

type ConsumerJobMarkCancelledRes struct {
	Updated map[string]bool
}

func (req *ConsumerJobMarkCancelledReq) Do(ctx context.Context, tx pgx.Tx) (*ConsumerJobMarkCancelledRes, error) {
	ctx, span := telemetry.Tracer().Start(ctx, "api_ConsumerJobMarkCancelled", trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()

	res := &ConsumerJobMarkCancelledRes{Updated: make(map[string]bool)}

	var names = make([]string, len(req.Events))
	var args = pgx.NamedArgs{
		"cancelled_state": int(entities.StateCancelled),
		"avaiable_state":  int(entities.StateAvailable),
		"running_state":   int(entities.StateRunning),
		"retryable_state": int(entities.StateRetryable),
		"finalized_at":    time.Now().UnixMilli(),
	}
	for i, event := range req.Events {
		// we assume that events are not able to update firstly
		// later if we can update its state, we can set it back to true
		res.Updated[event.EventId] = false

		bind := fmt.Sprintf("event_id_%d", i)
		names[i] = fmt.Sprintf("@%s", bind)
		args[bind] = event.EventId
	}

	table := pgx.Identifier{entities.CollectionConsumerJob(req.Consumer.Name)}.Sanitize()
	query := fmt.Sprintf(ConsumerJobMarkCancelledSQL, table, strings.Join(names, ","))

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			span.RecordError(err)
			return nil, err
		}

		res.Updated[id] = true
	}
	if rows.Err() != nil {
		span.RecordError(err)
		return nil, rows.Err()
	}

	var updated = make([]string, 0)
	var excluded = make([]string, 0)
	for id, ok := range res.Updated {
		if ok {
			updated = append(updated, id)
			continue
		}
		excluded = append(excluded, id)
	}
	span.SetAttributes(attribute.StringSlice("api_ConsumerJobMarkCancelled/updated", updated))
	span.SetAttributes(attribute.StringSlice("api_ConsumerJobMarkCancelled/excluded", excluded))

	return res, nil
}
