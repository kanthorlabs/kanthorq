package api

import (
	"context"
	"testing"

	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/stretchr/testify/require"
)

func TestConsumerJobMarkComplete(t *testing.T) {
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

		c, err := ConsumerEnsure(
			s.Stream,
			testify.ConsumerName(5),
			testify.Topic(5),
		).Do(ctx, tx)
		require.NoError(t, err)
		require.NotNil(t, c)
		require.NotNil(t, c.Consumer)

		// need to commit to make sure the consumer exists in next transaction
		require.NoError(t, tx.Commit(ctx))

		tx, err = pool.Begin(ctx)
		require.NoError(t, err)

		events := testify.GenStreamEvents(ctx, testify.Topic(5), 10)
		ids := make([]string, len(events))
		for i, event := range events {
			ids[i] = event.EventId
		}

		r, err := ConsumerJobMarkComplete(c.Consumer, ids).Do(ctx, tx)
		require.NoError(t, err)
		require.Equal(t, len(events), len(r.Updated))

		require.NoError(t, tx.Commit(ctx))
	})
}