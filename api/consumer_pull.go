package api

import (
	"context"
	"errors"
	"fmt"

	_ "embed"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
)

func ConsumerPull(consumer *entities.Consumer, size int) *ConsumerPullReq {
	return &ConsumerPullReq{
		Consumer: consumer,
		Size:     size,
	}
}

//go:embed consumer_pull.sql
var ConsumerPullSQL string

type ConsumerPullReq struct {
	Consumer *entities.Consumer
	Size     int
}

type ConsumerPullRes struct {
	CurrentCursor string
	NextCursor    string
}

func (req *ConsumerPullReq) Do(ctx context.Context, tx pgx.Tx) (*ConsumerPullRes, error) {
	cur, err := req.cursor().Do(ctx, tx)
	if err != nil {
		return nil, err
	}

	args := pgx.NamedArgs{
		"consumer_name":   req.Consumer.Name,
		"consumer_topic":  req.Consumer.Topic,
		"size":            req.Size,
		"consumer_cursor": cur.Cursor,
	}

	jtable := pgx.Identifier{entities.CollectionConsumerJob(req.Consumer.Name)}.Sanitize()
	stable := pgx.Identifier{entities.CollectionStreamEvent(req.Consumer.StreamName)}.Sanitize()
	query := fmt.Sprintf(ConsumerPullSQL, jtable, stable)

	var next string
	err = tx.QueryRow(ctx, query, args).Scan(&next)

	// in the initial state, there is no rows to pull
	// so pgx will return ErrNoRows, need to cast it as successful case
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	return &ConsumerPullRes{CurrentCursor: cur.Cursor, NextCursor: next}, nil
}

func (req *ConsumerPullReq) cursor() *ConsumerCursorReadReq {
	return &ConsumerCursorReadReq{
		Consumer: req.Consumer,
	}
}
