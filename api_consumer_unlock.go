package kanthorq

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/pkg/validator"
)

//go:embed api_consumer_unlock.sql
var ConsumerUnlockSql string

type ConsumerUnlockReq struct {
	Name   string `validate:"required"`
	Cursor string `validate:"required"`
}

type ConsumerUnlockRes struct {
	Consumer *ConsumerRegistry
}

func (req *ConsumerUnlockReq) Do(ctx context.Context, tx pgx.Tx) (*ConsumerUnlockRes, error) {
	err := validator.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	var res = &ConsumerUnlockRes{Consumer: &ConsumerRegistry{}}
	var args = pgx.NamedArgs{
		"consumer_name":   req.Name,
		"consumer_cursor": req.Cursor,
	}
	err = tx.QueryRow(ctx, ConsumerUnlockSql, args).Scan(
		&res.Consumer.StreamId,
		&res.Consumer.StreamName,
		&res.Consumer.Id,
		&res.Consumer.Name,
		&res.Consumer.SubjectFilter,
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
