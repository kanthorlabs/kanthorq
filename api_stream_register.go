package kanthorq

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/utils"
)

//go:embed api_stream_register_registry.sql
var StreamRegisterRegistrySql string

//go:embed api_stream_register_collection.sql
var StreamRegisterCollectionSql string

type StreamRegisterReq struct {
	StreamName string
}

type StreamRegisterRes struct {
	*StreamRegistry
}

func (req *StreamRegisterReq) Do(ctx context.Context, tx pgx.Tx) (*StreamRegisterRes, error) {
	// we are not sure we have the stream yet, so we cannot use row lock
	// must use advisory lock instead
	lock := utils.AdvisoryLockHash(req.StreamName)
	_, err := tx.Exec(ctx, fmt.Sprintf("SELECT pg_advisory_xact_lock(%d);", lock))
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
	table := pgx.Identifier{StreamCollection(req.StreamName)}.Sanitize()
	query := fmt.Sprintf(StreamRegisterCollectionSql, table, table)
	_, err = tx.Exec(ctx, query)

	return &StreamRegisterRes{&stream}, err
}
