package core

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/validator"
)

//go:embed stream_register_registry.sql
var StreamRegisterRegistrySql string

//go:embed stream_register_collection.sql
var StreamRegisterCollectionSql string

type StreamRegisterReq struct {
	StreamName string `validate:"required,is_collection_name"`
}

type StreamRegisterRes struct {
	*entities.StreamRegistry
}

func (req *StreamRegisterReq) Do(ctx context.Context, tx pgx.Tx) (*StreamRegisterRes, error) {
	err := validator.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	// register stream in registry
	var registry entities.StreamRegistry
	var args = pgx.NamedArgs{
		"stream_id":   entities.StreamId(),
		"stream_name": req.StreamName,
	}
	err = tx.
		QueryRow(ctx, StreamRegisterRegistrySql, args).
		Scan(&registry.Id, &registry.Name, &registry.CreatedAt, &registry.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// register stream collection
	table := pgx.Identifier{entities.Collection(registry.Id)}.Sanitize()
	query := fmt.Sprintf(StreamRegisterCollectionSql, table)
	_, err = tx.Exec(ctx, query)

	return &StreamRegisterRes{&registry}, err
}
