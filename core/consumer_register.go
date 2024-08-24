package core

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/xvalidator"
)

//go:embed consumer_register_registry.sql
var ConsumerRegisterRegistrySql string

//go:embed consumer_register_collection.sql
var ConsumerRegisterCollectionSql string

type ConsumerRegisterReq struct {
	StreamName            string   `validate:"required,is_collection_name"`
	ConsumerName          string   `validate:"required,is_collection_name"`
	ConsumerSubjectFilter []string `validate:"required,gt=0,dive,is_subject_filter"`
	ConsumerAttemptMax    int16    `validate:"required,gt=0"`
}

type ConsumerRegisterRes struct {
	*entities.StreamRegistry
	*entities.ConsumerRegistry
}

func (req *ConsumerRegisterReq) Do(ctx context.Context, tx pgx.Tx) (*ConsumerRegisterRes, error) {
	err := xvalidator.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	stream, err := (&StreamRegisterReq{StreamName: req.StreamName}).Do(ctx, tx)
	if err != nil {
		return nil, err
	}

	var registry entities.ConsumerRegistry
	var args = pgx.NamedArgs{
		"stream_id":               stream.StreamRegistry.Id,
		"stream_name":             stream.StreamRegistry.Name,
		"consumer_id":             entities.ConsumerId(),
		"consumer_name":           req.ConsumerName,
		"consumer_subject_filter": req.ConsumerSubjectFilter,
		"consumer_cursor":         entities.EventIdFromTime(time.UnixMilli(stream.StreamRegistry.CreatedAt)),
		"consumer_attempt_max":    req.ConsumerAttemptMax,
	}
	err = tx.
		QueryRow(ctx, ConsumerRegisterRegistrySql, args).
		Scan(
			&registry.StreamId,
			&registry.StreamName,
			&registry.Id,
			&registry.Name,
			&registry.SubjectFilter,
			&registry.Cursor,
			&registry.AttemptMax,
			&registry.CreatedAt,
			&registry.UpdatedAt,
		)
	if err != nil {
		return nil, err
	}

	// register stream collection
	table := pgx.Identifier{entities.Collection(registry.Id)}.Sanitize()
	query := fmt.Sprintf(ConsumerRegisterCollectionSql, table, table)
	_, err = tx.Exec(ctx, query)

	return &ConsumerRegisterRes{stream.StreamRegistry, &registry}, err
}