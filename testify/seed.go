package testify

import (
	"context"
	"fmt"
	"strings"
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

func GenStreamEvents(ctx context.Context, topic string, count int64) []*entities.StreamEvent {
	events := make([]*entities.StreamEvent, count)

	for i := 0; i < int(count); i++ {
		var now = time.Now().UTC().UnixMilli()
		events[i] = &entities.StreamEvent{
			Topic:     topic,
			EventId:   idx.New("evt"),
			CreatedAt: now,
		}
	}

	return events
}

func Topic(length int) string {
	words := Fake.Lorem().Words(length)

	var segments = make([]string, length)
	for i, word := range words {
		segments[i] = strings.ToLower(word)
	}
	return strings.Join(segments, ".")
}

func StreamName(length int) string {
	words := Fake.Lorem().Words(length)
	return strings.Join(words, "_")
}

func ConsumerName(length int) string {
	words := Fake.Lorem().Words(length)
	return strings.Join(words, "_")
}
