package kanthorq

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/pkg/validator"
)

//go:embed api_stream_register_registry.sql
var StreamRegisterRegistrySql string

//go:embed api_stream_register_collection.sql
var StreamRegisterCollectionSql string

type StreamRegisterReq struct {
	StreamName string `validate:"required,is_collection_name"`
}

type StreamRegisterRes struct {
	*StreamRegistry
}

func (req *StreamRegisterReq) Do(ctx context.Context, tx pgx.Tx) (*StreamRegisterRes, error) {
	err := validator.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	// register stream in registry
	var stream StreamRegistry
	var args = pgx.NamedArgs{"stream_name": req.StreamName}
	err = tx.
		QueryRow(ctx, StreamRegisterRegistrySql, args).
		Scan(&stream.Name, &stream.CreatedAt, &stream.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// register stream collection
	table := pgx.Identifier{Collection(req.StreamName)}.Sanitize()
	query := fmt.Sprintf(StreamRegisterCollectionSql, table, table)
	_, err = tx.Exec(ctx, query)

	return &StreamRegisterRes{&stream}, err
}
