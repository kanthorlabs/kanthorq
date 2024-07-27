package kanthorq

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/pkg/validator"
)

func ConsumerRegister(ctx context.Context, req *ConsumerRegisterReq, conn *pgx.Conn) (*ConsumerRegisterRes, error) {
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	res, err := req.Do(ctx, tx)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return res, nil
}

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

	var consumer ConsumerRegistry
	var args = pgx.NamedArgs{
		"stream_name":          stream.Name,
		"consumer_name":        req.ConsumerName,
		"consumer_topic":       req.ConsumerTopic,
		"consumer_attempt_max": req.ConsumerAttemptMax,
	}
	err = tx.
		QueryRow(ctx, ConsumerRegisterRegistrySql, args).
		Scan(
			&consumer.Name,
			&consumer.StreamName,
			&consumer.Topic,
			&consumer.Cursor,
			&consumer.AttemptMax,
			&consumer.CreatedAt,
			&consumer.UpdatedAt,
		)
	if err != nil {
		return nil, err
	}

	// register stream collection
	table := pgx.Identifier{Collection(consumer.Name)}.Sanitize()
	query := fmt.Sprintf(ConsumerRegisterCollectionSql, table, table)
	_, err = tx.Exec(ctx, query)

	return &ConsumerRegisterRes{stream.StreamRegistry, &consumer}, err
}
