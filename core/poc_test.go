package core

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kanthorlabs/common/idx"
	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/stretchr/testify/require"
)

var connection = os.Getenv("TEST_DATABASE_URI")

func BenchmarkPOC_ConsumerPull_DifferentSize(b *testing.B) {
	conn, err := pgxpool.New(context.Background(), connection)
	require.NoError(b, err)
	defer conn.Close()

	consumer := prepareConsumer(b, conn)
	for i := 1; i < 10; i++ {
		size := ConsumerPullSize * i
		b.Run(fmt.Sprintf("size %d", size), func(bs *testing.B) {
			cursor, err := QueryConsumerPull(consumer, size)(context.Background(), conn)
			require.NoError(bs, err)
			require.NotNil(bs, cursor.Name)
		})
	}
}

func BenchmarkPOC_ConsumerPull_MultipleConsumerReadSameTopic(b *testing.B) {
	conn, err := pgxpool.New(context.Background(), connection)
	require.NoError(b, err)
	defer conn.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		consumer := prepareConsumer(b, conn)
		for pb.Next() {
			cursor, err := QueryConsumerPull(consumer, ConsumerPullSize)(context.Background(), conn)
			require.NoError(b, err)
			require.NotNil(b, cursor.Name)
		}
	})
}

func prepareConsumer(b *testing.B, conn *pgxpool.Pool) string {
	var err error
	var consumer = idx.New("c")

	// truncate old jobs
	require.NoError(b, testify.QueryTruncateConsumer()(context.Background(), conn))
	// insert a fresh consumer
	_, err = QueryConsumerEnsure(consumer, os.Getenv("TEST_STREAM"), os.Getenv("TEST_TOPIC"))(context.Background(), conn)
	require.NoError(b, err)

	return consumer
}
