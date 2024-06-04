package publisher

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/api"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/q"
)

var _ Publisher = (*publisher)(nil)

func New(ctx context.Context, conf *Config) (Publisher, error) {
	conn, err := pgx.Connect(ctx, conf.ConnectionUri)
	if err != nil {
		return nil, err
	}

	stream, err := q.Stream(ctx, conn, &entities.Stream{
		Name: conf.StreamName,
	})
	if err != nil {
		return nil, err
	}

	return &publisher{
		conf:   conf,
		conn:   conn,
		stream: stream,
	}, nil
}

type publisher struct {
	conf *Config

	conn   *pgx.Conn
	stream *entities.Stream
}

func (pub *publisher) Send(ctx context.Context, events []*entities.StreamEvent) error {
	tx, err := pub.conn.Begin(ctx)
	if err != nil {
		return err
	}

	_, err = api.StreamEventPush(pub.stream, events).Do(ctx, tx)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
