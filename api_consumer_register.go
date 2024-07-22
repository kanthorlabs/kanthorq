package kanthorq

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5"
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
	StreamName         string
	ConsumerName       string
	ConsumerTopic      string
	ConsumerAttemptMax int16
}

type ConsumerRegisterRes struct {
	*StreamRegistry
	*ConsumerRegistry
}

func (req *ConsumerRegisterReq) Do(ctx context.Context, tx pgx.Tx) (*ConsumerRegisterRes, error) {
	stream, err := (&StreamRegisterReq{StreamName: req.StreamName}).Do(ctx, tx)
	if err != nil {
		return nil, err
	}

	var consumer ConsumerRegistry
	var args = pgx.NamedArgs{
		"stream_name":          req.StreamName,
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
	table := pgx.Identifier{ConsumerCollection(req.StreamName)}.Sanitize()
	query := fmt.Sprintf(ConsumerRegisterCollectionSql, table, table)
	_, err = tx.Exec(ctx, query)

	return &ConsumerRegisterRes{stream.StreamRegistry, &consumer}, err
}
