package q

import (
	"context"
	"testing"

	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/stretchr/testify/require"
)

func TestNewConsumerJobMarkRetryable(t *testing.T) {
	t.Run("happy case", func(t *testing.T) {
		ctx := context.Background()

		pool, err := testify.SetupPostgres(ctx)
		require.NoError(t, err)

		tx, err := pool.Begin(ctx)
		require.NoError(t, err)

		s, err := NewStreamEnsure(testify.StreamName(5)).Do(ctx, tx)
		require.NoError(t, err)
		require.NotNil(t, s)
		require.NotNil(t, s.Stream)

		c, err := NewConsumerEnsure(
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

		events := testify.GenStreamEvents(testify.Topic(5), 10)
		r, err := NewConsumerJobMarkRetryable(c.Consumer, events).Do(ctx, tx)
		require.NoError(t, err)
		require.Equal(t, len(events), len(r.Updated))

		require.NoError(t, tx.Commit(ctx))
	})
}
