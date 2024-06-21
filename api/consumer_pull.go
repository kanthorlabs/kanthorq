package api

import (
	"context"
	"errors"
	"fmt"

	_ "embed"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func NewConsumerPull(consumer *entities.Consumer, size int) *ConsumerPullReq {
	return &ConsumerPullReq{
		Consumer: consumer,
		Size:     size,
	}
}

//go:embed consumer_pull.sql
var ConsumerPullSQL string

type ConsumerPullReq struct {
	Consumer *entities.Consumer
	Size     int
}

type ConsumerPullRes struct {
	CurrentCursor string
	NextCursor    string
}

func (req *ConsumerPullReq) Do(ctx context.Context, tx pgx.Tx) (*ConsumerPullRes, error) {
	ctx, span := telemetry.Tracer().Start(ctx, "api_consumerpull", trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()

	cur, err := NewConsumerCursorRead(req.Consumer).Do(ctx, tx)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(attribute.String("consumer_name", req.Consumer.Name))
	span.SetAttributes(attribute.String("consumer_topic", req.Consumer.Topic))
	span.SetAttributes(attribute.Int("size", req.Size))
	span.SetAttributes(attribute.String("consumer_cursor", cur.Cursor))
	args := pgx.NamedArgs{
		"consumer_name":   req.Consumer.Name,
		"consumer_topic":  req.Consumer.Topic,
		"size":            req.Size,
		"consumer_cursor": cur.Cursor,
	}

	jtable := pgx.Identifier{entities.CollectionConsumerJob(req.Consumer.Name)}.Sanitize()
	stable := pgx.Identifier{entities.CollectionStreamEvent(req.Consumer.StreamName)}.Sanitize()
	query := fmt.Sprintf(ConsumerPullSQL, jtable, stable)

	res := &ConsumerPullRes{CurrentCursor: cur.Cursor, NextCursor: ""}
	err = tx.QueryRow(ctx, query, args).Scan(&res.NextCursor)

	// if error is nil or error is pgx.ErrNoRows, we accept it
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		span.RecordError(err)
		return nil, err
	}
	// if error is nil or error is pgx.ErrNoRows, we accept it

	return res, nil
}
