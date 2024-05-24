package core

import (
	"context"
	"fmt"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kanthorlabs/common/containers"
	"github.com/kanthorlabs/common/idx"
	"github.com/kanthorlabs/kanthorq/queries"
	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/stretchr/testify/require"
)

func BenchmarkPOC_ConsumerPull_DifferentSize(b *testing.B) {
	ctx := context.Background()
	container, err := testify.SpinPostgres(ctx, "kanthorlabs-kanthorq-storage-consumer-pull-size")
	require.NoError(b, err)

	uri, err := containers.PostgresConnectionString(ctx, container)
	require.NoError(b, err)

	conn, err := pgxpool.New(ctx, uri)
	require.NoError(b, err)
	defer conn.Close()

	stream := idx.New("stream")
	topic := idx.New("topic")
	if _, err := queries.EnsureStream(stream)(ctx, conn); err != nil {
		panic(err)
	}
	if err := testify.SeedStreamEvents(ctx, conn, stream, topic, 1000000); err != nil {
		panic(err)
	}

	consumer := idx.New("consumer")
	if _, err := queries.EnsureConsumer(consumer, stream, topic)(ctx, conn); err != nil {
		panic(err)
	}

	b.ResetTimer()
	for i := 0; i < 10; i++ {
		size := ConsumerPullSize * (i + 1)
		b.Run(fmt.Sprintf("size %d", size), func(bs *testing.B) {
			cursor, err := queries.ConsumerPull(consumer, size)(ctx, conn)
			require.NoError(bs, err)
			require.NotEmpty(bs, cursor.Name)
			require.NotEmpty(bs, cursor.Next)
		})
	}
}

func BenchmarkPOC_ConsumerPull_MultipleConsumerReadSameTopic(b *testing.B) {
	ctx := context.Background()
	container, err := testify.SpinPostgres(ctx, "kanthorlabs-kanthorq-storage-consumer-pull-concurrency")
	require.NoError(b, err)

	uri, err := containers.PostgresConnectionString(ctx, container)
	require.NoError(b, err)

	conn, err := pgxpool.New(ctx, uri)
	require.NoError(b, err)
	defer conn.Close()

	stream := idx.New("stream")
	topic := idx.New("topic")
	if _, err := queries.EnsureStream(stream)(ctx, conn); err != nil {
		panic(err)
	}
	if err := testify.SeedStreamEvents(ctx, conn, stream, topic, 1000000); err != nil {
		panic(err)
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			consumer := idx.New("c")
			_, err := queries.EnsureConsumer(consumer, stream, topic)(ctx, conn)
			require.NoError(b, err)

			cursor, err := queries.ConsumerPull(consumer, ConsumerPullSize)(ctx, conn)
			require.NoError(b, err)
			require.NotEmpty(b, cursor.Name)
			require.NotEmpty(b, cursor.Next)
		}
	})
}
