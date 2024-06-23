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

// NewConsumerJobMarkRetry mark events to be retried
// the reschedule time will be calculated based on the the formula:
// next_schedule_at = NOW() + ((attempt_count ^ 4) + (attempt_count ^ 4) * (RANDOM() * 0.2 - 0.1)) * 60 * 1000
// so there is the list of retries:
// First retry: ~ 1min
// Second retry: ~ 16mins
// Third retry: ~ 81mins = 1h21mins
func NewConsumerJobMarkRetry(consumer *entities.Consumer, events []*entities.StreamEvent) *ConsumerJobMarkRetryReq {
	return &ConsumerJobMarkRetryReq{
		Consumer: consumer,
		Events:   events,
	}
}

//go:embed consumer_job_mark_retry.sql
var ConsumerJobMarkRetrySQL string

type ConsumerJobMarkRetryReq struct {
	Consumer *entities.Consumer
	Events   []*entities.StreamEvent
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
		"consumer_name":   req.Consumer.Name,
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
	span.SetAttributes(attribute.StringSlice("api_consumerjobmarkcomplete_updated", updated))
	span.SetAttributes(attribute.StringSlice("api_consumerjobmarkcomplete_excluded", excluded))

	return res, nil
}
