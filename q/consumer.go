package q

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
)

func NewConsumer(ctx context.Context, conn *pgx.Conn, consumer *entities.Consumer) (*entities.Consumer, error) {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}

	s, err := NewStreamEnsure(consumer.StreamName).Do(ctx, tx)
	if err != nil {
		return nil, err
	}

	c, err := NewConsumerEnsure(s.Stream, consumer.Name, consumer.Topic).Do(ctx, tx)
	if err != nil {
		return nil, err
	}

	return c.Consumer, tx.Commit(ctx)
}
