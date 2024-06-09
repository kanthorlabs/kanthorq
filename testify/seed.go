package testify

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
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

	var size = 64 * 1024
	s, err := strconv.ParseInt(os.Getenv("KANTHORQ_EVENT_SIZE"), 10, 64)
	if err == nil && s > 0 {
		size = int(s)
	}

	for i := 0; i < int(count); i++ {
		var now = time.Now().UTC().UnixMilli()
		events[i] = &entities.StreamEvent{
			EventId:   ulid.Make().String(),
			Topic:     topic,
			Body:      GenBytes(size),
			Metadata:  GenBytes(size),
			CreatedAt: now,
		}
	}

	return events
}

func GenBytes(size int) []byte {
	if size <= 0 {
		panic("size must be greater than 0")
	}

	// Generate a random string of the given size
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	s := make([]rune, size)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}

	bytes, err := json.Marshal(string(s))
	if err != nil {
		panic(err)
	}

	return bytes
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
