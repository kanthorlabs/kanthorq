package kanthorq

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/utils"
)

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

type ConsumerRegisgerRes struct {
	*StreamRegistry
	*ConsumerRegistry
}

func (req *ConsumerRegisterReq) Do(ctx context.Context, tx pgx.Tx) (*ConsumerRegisgerRes, error) {
	stream, err := (&StreamRegisterReq{StreamName: req.StreamName}).Do(ctx, tx)
	if err != nil {
		return nil, err
	}

	// we are not sure we have the consumer yet, so we cannot use row lock
	// must use advisory lock instead
	lock := utils.AdvisoryLockHash(req.ConsumerName)
	_, err = tx.Exec(ctx, fmt.Sprintf("SELECT pg_advisory_xact_lock(%d);", lock))
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

	return &ConsumerRegisgerRes{stream.StreamRegistry, &consumer}, err
}
