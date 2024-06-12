package api

import (
	"context"
	"fmt"

	_ "embed"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/utils"
)

func ConsumerCreate(name string) *ConsumerCreateReq {
	return &ConsumerCreateReq{Name: name}
}

//go:embed consumer_create.sql
var ConsumerCreateSQL string

type ConsumerCreateReq struct {
	Name string
}

func (req *ConsumerCreateReq) Do(ctx context.Context, tx pgx.Tx) error {
	name := entities.CollectionConsumerJob(req.Name)
	table := pgx.Identifier{name}.Sanitize()
	query := fmt.Sprintf(ConsumerCreateSQL, utils.AdvisoryLockHash(name), table, table)

	_, err := tx.Exec(ctx, query)
	return err
}
