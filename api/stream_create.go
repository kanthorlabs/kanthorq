package api

import (
	"context"
	"fmt"

	_ "embed"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
)

func StreamCreate(name string) *StreamCreateReq {
	return &StreamCreateReq{Name: name}
}

//go:embed stream_create.sql
var StreamCreateSQL string

type StreamCreateReq struct {
	Name string
}

func (req *StreamCreateReq) Do(ctx context.Context, tx pgx.Tx) error {
	table := pgx.Identifier{entities.CollectionStreamEvent(req.Name)}.Sanitize()
	query := fmt.Sprintf(StreamCreateSQL, table, table)

	_, err := tx.Exec(ctx, query)
	return err
}