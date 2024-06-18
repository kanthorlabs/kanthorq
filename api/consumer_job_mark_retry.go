package api

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func NewConsumerJobMarkRetry(consumer *entities.Consumer, events []*entities.StreamEvent) *ConsumerJobMarkRetryReq {
	return &ConsumerJobMarkRetryReq{
		Consumer:   consumer,
		Events:     events,
		AttemptMax: 3,
	}
}

//go:embed consumer_job_mark_retry.sql
var ConsumerJobMarkRetrySQL string

type ConsumerJobMarkRetryReq struct {
	Consumer   *entities.Consumer
	Events     []*entities.StreamEvent
	AttemptMax int
}

type ConsumerJobMarkRetryRes struct {
	Updated map[string]bool
}

func (req *ConsumerJobMarkRetryReq) Do(ctx context.Context, tx pgx.Tx) (*ConsumerJobMarkRetryRes, error) {
	ctx, span := telemetry.Tracer().Start(ctx, "api_consumerjobmarkretry", trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()

	res := &ConsumerJobMarkRetryRes{Updated: make(map[string]bool)}

	var names = make([]string, len(req.Events))
	var args = pgx.NamedArgs{
		"attempt_max":     req.AttemptMax,
		"discarded_state": int(entities.StateDiscarded),
		"retryable_state": int(entities.StateRetryable),
		"running_state":   int(entities.StateRunning),
	}
	for i, event := range req.Events {
		// we assume that events are not able to update firstly
		// later if we can update its state, we can set it back to true
		res.Updated[event.EventId] = false

		bind := fmt.Sprintf("event_id_%d", i)
		names[i] = fmt.Sprintf("@%s", fmt.Sprintf("event_id_%d", i))
		args[bind] = event.EventId
	}

	table := pgx.Identifier{entities.CollectionConsumerJob(req.Consumer.Name)}.Sanitize()
	query := fmt.Sprintf(ConsumerJobMarkRetrySQL, table, strings.Join(names, ","))

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
	span.SetAttributes(attribute.StringSlice("api.ConsumerJobMarkComplete/updated", updated))
	span.SetAttributes(attribute.StringSlice("api.ConsumerJobMarkComplete/excluded", excluded))

	return res, nil
}
