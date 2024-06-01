package api

import (
	"context"
	"testing"

	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/stretchr/testify/require"
)

func TestConsumerPull(t *testing.T) {
	t.Run("happy case", func(t *testing.T) {
		ctx := context.Background()

		pool, err := testify.SetupPostgres(ctx)
		require.NoError(t, err)

		tx, err := pool.Begin(ctx)
		require.NoError(t, err)

		s, err := StreamEnsure(testify.Fake.RandomStringWithLength(32)).Do(ctx, tx)
		require.NoError(t, err)
		require.NotNil(t, s)
		require.NotNil(t, s.Stream)

		c, err := ConsumerEnsure(
			s.Stream,
			testify.Fake.RandomStringWithLength(32),
			testify.Fake.RandomStringWithLength(32),
		).Do(ctx, tx)
		require.NoError(t, err)
		require.NotNil(t, c)
		require.NotNil(t, c.Consumer)

		// need to commit to make sure the consumer exists in next transaction
		require.NoError(t, tx.Commit(ctx))

		tx, err = pool.Begin(ctx)
		require.NoError(t, err)

		changes, err := ConsumerPull(c.Consumer, testify.Fake.IntBetween(10, 100)).Do(ctx, tx)
		require.NoError(t, err)
		require.NotNil(t, changes)

		require.NoError(t, tx.Commit(ctx))

	})
}
