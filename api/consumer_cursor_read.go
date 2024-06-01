package api

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
)

func ConsumerCursorRead(Consumer *entities.Consumer) *ConsumerCursorReadReq {
	return &ConsumerCursorReadReq{Consumer}
}

//go:embed consumer_cursor_read.sql
var ConsumerCursorReadSQL string

type ConsumerCursorReadReq struct {
	Consumer *entities.Consumer
}

type ConsumerCursorReadRes struct {
	Cursor string
}

func (req *ConsumerCursorReadReq) Do(ctx context.Context, tx pgx.Tx) (*ConsumerCursorReadRes, error) {
	args := pgx.NamedArgs{"consumer_name": req.Consumer.Name}

	table := pgx.Identifier{entities.CollectionConsumer}.Sanitize()
	query := fmt.Sprintf(ConsumerCursorReadSQL, table)

	var cursor *string
	if err := tx.QueryRow(ctx, query, args).Scan(&cursor); err != nil {
		return nil, err
	}
	if cursor == nil {
		return nil, fmt.Errorf("ERROR.CONSUMER.BUSY: %s", req.Consumer.Name)
	}

	return &ConsumerCursorReadRes{Cursor: *cursor}, nil
}
