package testify

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kanthorlabs/common/idx"
	"github.com/kanthorlabs/kanthorq/entities"
)

func SeedStreamEvents(ctx context.Context, conn *pgxpool.Pool, stream, topic string, count int) error {
	events := make([][]any, count)
	for i := 0; i < count; i++ {
		var now = time.Now().UTC().UnixMilli()
		events[i] = []any{topic, idx.New("evt"), fmt.Sprintf("%d", now)}
	}

	_, err := conn.CopyFrom(
		ctx,
		pgx.Identifier{entities.CollectionStreamEvent(stream)},
		entities.CollectionStreamEventProps,
		pgx.CopyFromRows(events),
	)
	return err
}
