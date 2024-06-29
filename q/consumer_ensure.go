package q

import (
	"context"

	_ "embed"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
)

func NewConsumerEnsure(stream *entities.Stream, name, topic string) *ConsumerEnsureReq {
	return &ConsumerEnsureReq{Stream: stream, Name: name, Topic: topic}
}

//go:embed consumer_ensure.sql
var ConsumerEnsureSQL string

type ConsumerEnsureReq struct {
	Stream *entities.Stream
	Name   string
	Topic  string
}
type ConsumerEnsureRes struct {
	*entities.Consumer
}

func (req *ConsumerEnsureReq) Do(ctx context.Context, tx pgx.Tx) (*ConsumerEnsureRes, error) {
	args := pgx.NamedArgs{
		"stream_name":   req.Stream.Name,
		"consumer_name": req.Name,
		"topic":         req.Topic,
	}

	var consumer entities.Consumer
	err := tx.
		QueryRow(ctx, ConsumerEnsureSQL, args).
		Scan(
			&consumer.Name,
			&consumer.StreamName,
			&consumer.Topic,
			&consumer.Cursor,
			&consumer.CreatedAt,
			&consumer.UpdatedAt,
		)
	if err != nil {
		return nil, err
	}

	if err := NewConsumerCreate(consumer.Name).Do(ctx, tx); err != nil {
		return nil, err
	}

	return &ConsumerEnsureRes{&consumer}, nil
}
