package kanthorq

import (
	"context"
	"errors"
	"fmt"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kanthorlabs/common/idx"
	"github.com/kanthorlabs/kanthorq/queries"
	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/stretchr/testify/require"
)

func Benchmark_ConsumerPull_DifferentSize(b *testing.B) {
	ctx := context.Background()
	conn, err := testify.SetupPostgres(ctx)
	require.NoError(b, err)
	defer conn.Close()

	stream := idx.New("s")
	topic := idx.New("topic")

	tx, err := conn.Begin(ctx)
	require.NoError(b, err)
	if _, err := queries.EnsureStream(stream)(ctx, tx); err != nil {
		panic(errors.Join(err, tx.Rollback(ctx)))
	}
	consumer, err := queries.EnsureConsumer(idx.New("job"), stream, topic)(ctx, tx)
	if err != nil {
		panic(errors.Join(err, tx.Rollback(ctx)))
	}
	require.NoError(b, tx.Commit(ctx))

	if err := testify.SeedStreamEvents(ctx, conn, stream, topic, 1000000); err != nil {
		panic(err)
	}

	b.ResetTimer()
	for i := 0; i < 10; i++ {
		size := ConsumerPullSize * (i + 1)
		b.Run(fmt.Sprintf("size %d", size), func(bs *testing.B) {
			tx, err := conn.Begin(ctx)
			require.NoError(bs, err)
			cursor, err := queries.ConsumerPull(consumer, size)(ctx, tx)
			require.NoError(bs, err)
			require.NotEmpty(bs, cursor)
			require.NoError(bs, tx.Commit(ctx))
		})
	}
}

func Benchmark_ConsumerPull_MultipleConsumerReadSameTopic(b *testing.B) {
	ctx := context.Background()
	conn, err := testify.SetupPostgres(ctx)
	require.NoError(b, err)
	defer conn.Close()

	stream := idx.New("s")
	topic := idx.New("topic")

	tx, err := conn.Begin(ctx)
	require.NoError(b, err)
	if _, err := queries.EnsureStream(stream)(ctx, tx); err != nil {
		panic(errors.Join(err, tx.Rollback(ctx)))
	}
	require.NoError(b, tx.Commit(ctx))

	if err := testify.SeedStreamEvents(ctx, conn, stream, topic, 1000000); err != nil {
		panic(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tx, err := conn.Begin(ctx)
			require.NoError(b, err)

			consumer, err := queries.EnsureConsumer(idx.New("job"), stream, topic)(ctx, tx)
			require.NoError(b, err)

			cursor, err := queries.ConsumerPull(consumer, ConsumerPullSize)(ctx, tx)
			require.NoError(b, err)

			require.NotEmpty(b, cursor)
			require.NoError(b, tx.Commit(ctx))
		}
	})
}

func Benchmark_ConsumerJobPull_DifferentSize(b *testing.B) {
	ctx := context.Background()
	conn, err := testify.SetupPostgres(ctx)
	require.NoError(b, err)
	defer conn.Close()

	stream := idx.New("s")
	topic := idx.New("topic")

	tx, err := conn.Begin(ctx)
	require.NoError(b, err)
	if _, err := queries.EnsureStream(stream)(ctx, tx); err != nil {
		panic(errors.Join(err, tx.Rollback(ctx)))
	}
	consumer, err := queries.EnsureConsumer(idx.New("job"), stream, topic)(ctx, tx)
	if err != nil {
		panic(errors.Join(err, tx.Rollback(ctx)))
	}
	require.NoError(b, tx.Commit(ctx))

	if err := testify.SeedStreamEvents(ctx, conn, stream, topic, 1000000); err != nil {
		panic(err)
	}

	b.ResetTimer()
	for i := 0; i < 10; i++ {
		size := ConsumerJobPullSize * (i + 1)
		b.Run(fmt.Sprintf("size %d", size), func(bs *testing.B) {
			tx, err := conn.Begin(ctx)
			require.NoError(b, err)

			cursor, err := queries.ConsumerPull(consumer, size)(ctx, tx)
			require.NoError(b, err)
			require.NotEmpty(b, cursor)

			events, err := queries.ConsumerJobPull(consumer, size)(ctx, tx)

			require.NoError(bs, err)
			require.NotNil(bs, events)
			require.Equal(bs, size, len(events))
			require.NoError(b, tx.Commit(ctx))
		})
	}
}

func Benchmark_ConsumerJobPull_MultipleConsumerReadSameTopic(b *testing.B) {
	ctx := context.Background()
	conn, err := testify.SetupPostgres(ctx)
	require.NoError(b, err)
	defer conn.Close()

	stream := idx.New("s")
	topic := idx.New("topic")

	tx, err := conn.Begin(ctx)
	require.NoError(b, err)
	if _, err := queries.EnsureStream(stream)(ctx, tx); err != nil {
		panic(errors.Join(err, tx.Rollback(ctx)))
	}
	require.NoError(b, tx.Commit(ctx))

	if err := testify.SeedStreamEvents(ctx, conn, stream, topic, 1000000); err != nil {
		panic(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tx, err := conn.Begin(ctx)
			require.NoError(b, err)

			consumer, err := queries.EnsureConsumer(idx.New("job"), stream, topic)(ctx, tx)
			require.NoError(b, err)

			cursor, err := queries.ConsumerPull(consumer, ConsumerJobPullSize)(ctx, tx)
			require.NoError(b, err)
			require.NotEmpty(b, cursor)

			events, err := queries.ConsumerJobPull(consumer, ConsumerJobPullSize)(ctx, tx)

			require.NoError(b, err)
			require.NotNil(b, events)
			require.Equal(b, ConsumerJobPullSize, len(events))

			require.NoError(b, tx.Commit(ctx))
		}
	})
}
