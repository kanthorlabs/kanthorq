package kanthorq

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/pkg/validator"
)

//go:embed api_consumer_register_registry.sql
var ConsumerRegisterRegistrySql string

//go:embed api_consumer_register_collection.sql
var ConsumerRegisterCollectionSql string

type ConsumerRegisterReq struct {
	StreamName         string `validate:"required,is_collection_name"`
	ConsumerName       string `validate:"required,is_collection_name"`
	ConsumerTopic      string `validate:"required,is_topic"`
	ConsumerAttemptMax int16  `validate:"required,gt=0"`
}

type ConsumerRegisterRes struct {
	*StreamRegistry
	*ConsumerRegistry
}

func (req *ConsumerRegisterReq) Do(ctx context.Context, tx pgx.Tx) (*ConsumerRegisterRes, error) {
	err := validator.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	stream, err := (&StreamRegisterReq{StreamName: req.StreamName}).Do(ctx, tx)
	if err != nil {
		return nil, err
	}

	var registry ConsumerRegistry
	var args = pgx.NamedArgs{
		"stream_id":            stream.StreamRegistry.Id,
		"stream_name":          stream.StreamRegistry.Name,
		"consumer_id":          ConsumerId(),
		"consumer_name":        req.ConsumerName,
		"consumer_topic":       req.ConsumerTopic,
		"consumer_attempt_max": req.ConsumerAttemptMax,
	}
	err = tx.
		QueryRow(ctx, ConsumerRegisterRegistrySql, args).
		Scan(
			&registry.StreamId,
			&registry.StreamName,
			&registry.Id,
			&registry.Name,
			&registry.Topic,
			&registry.Cursor,
			&registry.AttemptMax,
			&registry.CreatedAt,
			&registry.UpdatedAt,
		)
	if err != nil {
		return nil, err
	}

	// register stream collection
	table := pgx.Identifier{Collection(registry.Id)}.Sanitize()
	query := fmt.Sprintf(ConsumerRegisterCollectionSql, table, table)
	_, err = tx.Exec(ctx, query)

	return &ConsumerRegisterRes{stream.StreamRegistry, &registry}, err
}
