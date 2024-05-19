package core

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/common/idx"
	"github.com/stretchr/testify/require"
)

func BenchmarkPOC_ConsumerPull_DifferentSize(b *testing.B) {
	var uri = os.Getenv("TEST_DATABASE_URI")

	consumer := prepareConsumer(b)
	for i := 1; i < 10; i++ {
		size := ConsumerPullSize * i
		b.Run(fmt.Sprintf("size %d", size), func(bs *testing.B) {
			conn, err := pgx.Connect(context.Background(), uri)
			require.NoError(bs, err)
			defer conn.Close(context.Background())

			var cursor string
			err = conn.
				QueryRow(
					context.Background(),
					QueryConsumerPull,
					pgx.NamedArgs{
						"consumer_name": consumer,
						"size":          size,
					},
				).
				Scan(&cursor)
			require.NoError(bs, err)
			require.NotEmpty(bs, cursor)
		})
	}
}

func BenchmarkPOC_ConsumerPull_MultipleConsumerReadSameTopic(b *testing.B) {
	var uri = os.Getenv("TEST_DATABASE_URI")
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		consumer := prepareConsumer(b)

		conn, err := pgx.Connect(context.Background(), uri)
		require.NoError(b, err)
		defer conn.Close(context.Background())

		var cursor string
		for pb.Next() {
			err = conn.
				QueryRow(
					context.Background(),
					QueryConsumerPull,
					pgx.NamedArgs{
						"consumer_name": consumer,
						"size":          ConsumerPullSize,
					},
				).
				Scan(&cursor)
			require.NoError(b, err)
			require.NotEmpty(b, cursor)
		}
	})
}

func prepareConsumer(b *testing.B) string {
	var uri = os.Getenv("TEST_DATABASE_URI")
	conn, err := pgx.Connect(context.Background(), uri)
	require.NoError(b, err)
	defer conn.Close(context.Background())

	var consumer = idx.New("cs")

	// truncate old jobs
	_, err = conn.Exec(context.Background(), QueryTruncate(CollectionConsumerJob))
	require.NoError(b, err)
	// delete old consumer
	_, err = conn.Exec(context.Background(), QueryTruncate(CollectionConsumer))
	require.NoError(b, err)
	// insert a fresh consumer
	_, err = conn.Exec(
		context.Background(),
		fmt.Sprintf(`INSERT INTO %s (name, tier, topic, cursor) VALUES (@name, @tier, @topic, '')`, CollectionConsumer),
		pgx.NamedArgs{
			"name":  consumer,
			"tier":  os.Getenv("TEST_TIER"),
			"topic": os.Getenv("TEST_TOPIC"),
		},
	)
	require.NoError(b, err)

	return consumer
}
