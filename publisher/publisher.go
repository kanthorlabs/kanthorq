package publisher

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/api"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/q"
)

var _ Publisher = (*publisher)(nil)

func New(conf *Config) Publisher {
	return &publisher{conf: conf}
}

type publisher struct {
	conf *Config

	conn   *pgx.Conn
	stream *entities.Stream
}

func (pub *publisher) Start(ctx context.Context) error {
	conn, err := pgx.Connect(ctx, pub.conf.ConnectionUri)
	if err != nil {
		return err
	}
	pub.conn = conn

	stream, err := q.Stream(ctx, pub.conn, &entities.Stream{Name: pub.conf.StreamName})
	if err != nil {
		return err
	}
	pub.stream = stream

	return nil
}

func (pub *publisher) Stop(ctx context.Context) error {
	return pub.conn.Close(ctx)
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
