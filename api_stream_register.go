package kanthorq

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/utils"
)

func StreamRegister(name string) *StreamRegisterReq {
	return &StreamRegisterReq{Name: name}
}

//go:embed api_stream_register_registry.sql
var StreamRegisterRegistrySql string

//go:embed api_stream_register_collection.sql
var StreamRegisterCollectionSql string

type StreamRegisterReq struct {
	Name string
}

type StreamRegisgerRes struct {
	*StreamRegistry
}

func (req *StreamRegisterReq) Do(ctx context.Context, tx pgx.Tx) (*StreamRegisgerRes, error) {
	lock := utils.AdvisoryLockHash(req.Name)
	_, err := tx.Exec(ctx, fmt.Sprintf("SELECT pg_advisory_xact_lock(%d);", lock))
	if err != nil {
		return nil, err
	}

	// register stream in registry
	var stream StreamRegistry
	var args = pgx.NamedArgs{"stream_name": req.Name}
	err = tx.
		QueryRow(ctx, StreamRegisterRegistrySql, args).
		Scan(&stream.Name, &stream.CreatedAt, &stream.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// register stream collection
	table := pgx.Identifier{StreamCollection(req.Name)}.Sanitize()
	query := fmt.Sprintf(StreamRegisterCollectionSql, table, table)
	_, err = tx.Exec(ctx, query)

	return &StreamRegisgerRes{&stream}, err
}
