package api

import (
	"context"

	_ "embed"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
)

func StreamEnsure(name string) *StreamEnsureReq {
	return &StreamEnsureReq{Name: name}
}

//go:embed stream_ensure.sql
var StreamEnsureSQL string

type StreamEnsureReq struct {
	Name string
}
type StreamEnsureRes struct {
	*entities.Stream
}

func (req *StreamEnsureReq) Do(ctx context.Context, tx pgx.Tx) (*StreamEnsureRes, error) {
	args := pgx.NamedArgs{"stream_name": req.Name}

	var stream entities.Stream
	err := tx.
		QueryRow(ctx, StreamEnsureSQL, args).
		Scan(&stream.Name, &stream.CreatedAt, &stream.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if err := StreamCreate(stream.Name).Do(ctx, tx); err != nil {
		return nil, err
	}

	return &StreamEnsureRes{&stream}, nil
}
