package core

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/xvalidator"
)

//go:embed consumer_unlock.sql
var ConsumerUnlockSql string

type ConsumerUnlockReq struct {
	Name   string `validate:"required"`
	Cursor string `validate:"required"`
}

type ConsumerUnlockRes struct {
	Consumer *entities.ConsumerRegistry
}

func (req *ConsumerUnlockReq) Do(ctx context.Context, tx pgx.Tx) (*ConsumerUnlockRes, error) {
	err := xvalidator.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	var res = &ConsumerUnlockRes{Consumer: &entities.ConsumerRegistry{}}
	var args = pgx.NamedArgs{
		"consumer_name":   req.Name,
		"consumer_cursor": req.Cursor,
	}
	err = tx.QueryRow(ctx, ConsumerUnlockSql, args).Scan(
		&res.Consumer.StreamId,
		&res.Consumer.StreamName,
		&res.Consumer.Id,
		&res.Consumer.Name,
		&res.Consumer.Kind,
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
