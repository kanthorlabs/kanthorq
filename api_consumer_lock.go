package kanthorq

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/pkg/validator"
)

//go:embed api_consumer_lock.sql
var ConsumerLockSql string

type ConsumerLockReq struct {
	Name string `validate:"required"`
}

type ConsumerLockRes struct {
	Consumer *ConsumerRegistry
}

func (req *ConsumerLockReq) Do(ctx context.Context, tx pgx.Tx) (*ConsumerLockRes, error) {
	err := validator.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	var res = &ConsumerLockRes{Consumer: &ConsumerRegistry{}}
	var args = pgx.NamedArgs{
		"consumer_name": req.Name,
	}
	err = tx.QueryRow(ctx, ConsumerLockSql, args).Scan(
		&res.Consumer.StreamId,
		&res.Consumer.StreamName,
		&res.Consumer.Id,
		&res.Consumer.Name,
		&res.Consumer.Subject,
		&res.Consumer.Cursor,
		&res.Consumer.AttemptMax,
		&res.Consumer.CreatedAt,
		&res.Consumer.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}
