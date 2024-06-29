package q

import (
	"context"
	"fmt"

	_ "embed"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/utils"
)

func NewConsumerJobCreate(name string) *ConsumerJobCreateReq {
	return &ConsumerJobCreateReq{Name: name}
}

//go:embed consumer_job_create.sql
var ConsumerJobCreateSQL string

type ConsumerJobCreateReq struct {
	Name string
}

func (req *ConsumerJobCreateReq) Do(ctx context.Context, tx pgx.Tx) error {
	lock := utils.AdvisoryLockHash(req.Name)
	table := pgx.Identifier{entities.CollectionConsumerJob(req.Name)}.Sanitize()
	query := fmt.Sprintf(ConsumerJobCreateSQL, lock, table, table)

	_, err := tx.Exec(ctx, query)
	return err
}
