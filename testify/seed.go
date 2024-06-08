package testify

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/oklog/ulid/v2"
)

func SeedStreamEvents(ctx context.Context, conn *pgxpool.Pool, stream, topic string, count int) error {
	events := make([][]any, count)
	for i := 0; i < count; i++ {
		var now = time.Now().UTC().UnixMilli()
		events[i] = []any{topic, ulid.Make().String(), fmt.Sprintf("%d", now)}
	}

	_, err := conn.CopyFrom(
		ctx,
		pgx.Identifier{entities.CollectionStreamEvent(stream)},
		entities.CollectionStreamEventProps,
		pgx.CopyFromRows(events),
	)
	return err
}

func GenStreamEvents(ctx context.Context, topic string, count int64) []*entities.StreamEvent {
	events := make([]*entities.StreamEvent, count)

	for i := 0; i < int(count); i++ {
		var now = time.Now().UTC().UnixMilli()
		events[i] = &entities.StreamEvent{
			Topic:     topic,
			EventId:   ulid.Make().String(),
			CreatedAt: now,
		}
	}

	return events
}

func Topic(wc int) string {
	words := Fake.Lorem().Words(wc)

	var segments = make([]string, wc)
	for i, word := range words {
		segments[i] = strings.ToLower(word)
	}
	return strings.Join(segments, ".")
}

func StreamName(wc int) string {
	words := Fake.Lorem().Words(wc)
	return strings.Join(words, "_")
}

func ConsumerName(wc int) string {
	words := Fake.Lorem().Words(wc)
	return strings.Join(words, "_")
}
