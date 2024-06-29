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

func NewConsumerJobMarkCompleted(consumer *entities.Consumer, events []*entities.StreamEvent) *ConsumerJobMarkCompletedReq {
	return &ConsumerJobMarkCompletedReq{
		Consumer: consumer,
		Events:   events,
	}
}

//go:embed consumer_job_mark_completed.sql
var ConsumerJobMarkCompletedSQL string

type ConsumerJobMarkCompletedReq struct {
	Consumer *entities.Consumer
	Events   []*entities.StreamEvent
}

type ConsumerJobMarkCompletedRes struct {
	Updated map[string]bool
}

func (req *ConsumerJobMarkCompletedReq) Do(ctx context.Context, tx pgx.Tx) (*ConsumerJobMarkCompletedRes, error) {
	ctx, span := telemetry.Tracer().Start(ctx, "api_ConsumerJobMarkCompleted", trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()

	res := &ConsumerJobMarkCompletedRes{Updated: make(map[string]bool)}

	var names = make([]string, len(req.Events))
	var args = pgx.NamedArgs{
		"running_state":   int(entities.StateRunning),
		"completed_state": int(entities.StateCompleted),
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
	query := fmt.Sprintf(ConsumerJobMarkCompletedSQL, table, strings.Join(names, ","))

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
	span.SetAttributes(attribute.StringSlice("api_ConsumerJobMarkCompleted/updated", updated))
	span.SetAttributes(attribute.StringSlice("api_ConsumerJobMarkCompleted/excluded", excluded))

	return res, nil
}
