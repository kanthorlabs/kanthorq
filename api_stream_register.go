package kanthorq

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func StreamRegister(ctx context.Context, req *StreamRegisterReq, conn *pgx.Conn) (*StreamRegisterRes, error) {
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	res, err := req.Do(ctx, tx)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return res, nil
}

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
	var err error

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
	query := fmt.Sprintf(StreamRegisterCollectionSql, table)
	_, err = tx.Exec(ctx, query)

	return &StreamRegisterRes{&stream}, err
}
