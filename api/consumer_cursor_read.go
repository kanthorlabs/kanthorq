package api

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func NewConsumerCursorRead(consumer *entities.Consumer) *ConsumerCursorReadReq {
	return &ConsumerCursorReadReq{Consumer: consumer}
}

//go:embed consumer_cursor_read.sql
var ConsumerCursorReadSQL string

type ConsumerCursorReadReq struct {
	Consumer *entities.Consumer
}

type ConsumerCursorReadRes struct {
	Cursor string
}

func (req *ConsumerCursorReadReq) Do(ctx context.Context, tx pgx.Tx) (*ConsumerCursorReadRes, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "api.ConsumerCursorRead", trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()

	span.SetAttributes(attribute.String("consumer_name", req.Consumer.Name))
	args := pgx.NamedArgs{"consumer_name": req.Consumer.Name}

	table := pgx.Identifier{entities.CollectionConsumer}.Sanitize()
	query := fmt.Sprintf(ConsumerCursorReadSQL, table)

	res := &ConsumerCursorReadRes{Cursor: ""}
	var err = tx.QueryRow(ctx, query, args).Scan(&res.Cursor)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		span.RecordError(err)
		return nil, err
	}
	// if error is nil or error is pgx.ErrNoRows, we accept it

	span.SetAttributes(attribute.String("consumer_cursor", res.Cursor))
	return res, nil
}
