package core

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/xvalidator"
)

//go:embed consumer_lock.sql
var ConsumerLockSql string

type ConsumerLockReq struct {
	Name string `validate:"required"`
}

type ConsumerLockRes struct {
	Consumer *entities.ConsumerRegistry
}

func (req *ConsumerLockReq) Do(ctx context.Context, tx pgx.Tx) (*ConsumerLockRes, error) {
	err := xvalidator.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	var res = &ConsumerLockRes{Consumer: &entities.ConsumerRegistry{}}
	var args = pgx.NamedArgs{
		"consumer_name": req.Name,
	}
	err = tx.QueryRow(ctx, ConsumerLockSql, args).Scan(
		&res.Consumer.StreamId,
		&res.Consumer.StreamName,
		&res.Consumer.Id,
		&res.Consumer.Name,
		&res.Consumer.SubjectIncludes,
		&res.Consumer.SubjectExcludes,
		&res.Consumer.Cursor,
		&res.Consumer.AttemptMax,
		&res.Consumer.VisibilityTimeout,
		&res.Consumer.CreatedAt,
		&res.Consumer.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}
