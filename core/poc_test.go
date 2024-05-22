package core

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/common/idx"
	"github.com/kanthorlabs/kanthorq/testify"
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

			var name, current, next *string
			statement, args := QueryConsumerPull(consumer, size)
			err = conn.
				QueryRow(context.Background(), statement, args).
				Scan(&name, &current, &next)

			require.NoError(bs, err)
			require.NotNil(bs, name)
			require.NotEmpty(bs, *name)
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

		for pb.Next() {

			var name, current, next *string
			statement, args := QueryConsumerPull(consumer, ConsumerPullSize)
			err = conn.
				QueryRow(context.Background(), statement, args).
				Scan(&name, &current, &next)

			require.NoError(b, err)
			require.NotNil(b, name)
			require.NotEmpty(b, *name)
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
	_, err = conn.Exec(context.Background(), testify.QueryTruncateConsumer())
	require.NoError(b, err)
	// insert a fresh consumer
	statement, args := QueryConsumerEnsure(consumer, os.Getenv("TEST_TOPIC"))
	_, err = conn.Exec(context.Background(), statement, args)
	require.NoError(b, err)

	return consumer
}
