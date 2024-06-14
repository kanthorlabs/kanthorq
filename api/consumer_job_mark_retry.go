package api

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

func ConsumerJobMarkRetry(consumer *entities.Consumer, eventIds []string) *ConsumerJobMarkRetryReq {
	return &ConsumerJobMarkRetryReq{
		Consumer: consumer,
		EventIds: eventIds,
	}
}

//go:embed consumer_job_mark_retry.sql
var ConsumerJobMarkRetrySQL string

type ConsumerJobMarkRetryReq struct {
	Consumer *entities.Consumer
	EventIds []string
}

type ConsumerJobMarkRetryRes struct {
	Updated map[string]bool
}

func (req *ConsumerJobMarkRetryReq) Do(ctx context.Context, tx pgx.Tx) (*ConsumerJobMarkRetryRes, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "api.ConsumerJobMarkRetry", trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()

	res := &ConsumerJobMarkRetryRes{Updated: make(map[string]bool)}

	var names = make([]string, len(req.EventIds))
	var args = pgx.NamedArgs{
		"retry_state":   entities.StateRetryable,
		"running_state": entities.StateRunning,
		"attempt_at":    time.Now().UnixMilli(),
	}
	for i, id := range req.EventIds {
		// we assume that events are not able to update firstly
		// later if we can update its state, we can set it back to true
		res.Updated[id] = false

		eventIdBind := fmt.Sprintf("event_id_%d", i)
		names[i] = fmt.Sprintf("@%s", eventIdBind)
		args[eventIdBind] = id
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
