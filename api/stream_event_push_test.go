package api

import (
	"context"
	"testing"

	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/stretchr/testify/require"
)

func TestStreamEventPush(t *testing.T) {
	t.Run("happy case", func(t *testing.T) {
		ctx := context.Background()

		pool, err := testify.SetupPostgres(ctx)
		require.NoError(t, err)

		tx, err := pool.Begin(ctx)
		require.NoError(t, err)

		s, err := StreamEnsure(testify.StreamName(5)).Do(ctx, tx)
		require.NoError(t, err)
		require.NotNil(t, s)
		require.NotNil(t, s.Stream)

		require.NoError(t, tx.Commit(ctx))

		tx, err = pool.Begin(ctx)
		require.NoError(t, err)

		count := testify.Fake.Int64Between(1000, 5000)
		events := testify.GenStreamEvents(ctx, testify.StreamName(5), count)
		p, err := StreamEventPush(s.Stream, events).Do(ctx, tx)
		require.NoError(t, err)
		require.NotNil(t, p)
		require.Equal(t, count, p.InsertCount)

		require.NoError(t, tx.Commit(ctx))

	})
}
