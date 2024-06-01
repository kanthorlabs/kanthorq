package api

import (
	"context"
	"testing"

	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/stretchr/testify/require"
)

func TestConsumerCursorRead(t *testing.T) {
	t.Run("happy case", func(t *testing.T) {
		ctx := context.Background()

		pool, err := testify.SetupPostgres(ctx)
		require.NoError(t, err)

		// seed an consumer
		tx, err := pool.Begin(ctx)
		require.NoError(t, err)

		ensure, err := ConsumerEnsure(
			testify.Fake.RandomStringWithLength(32),
			testify.Fake.RandomStringWithLength(32),
			testify.Fake.RandomStringWithLength(32),
		).Do(ctx, tx)
		require.NoError(t, err)
		require.NotNil(t, ensure)
		require.NotNil(t, ensure.Consumer)

		// need to commit to make sure the consumer exists in next transaction
		require.NoError(t, tx.Commit(ctx))

		tx, err = pool.Begin(ctx)
		require.NoError(t, err)

		res, err := ConsumerCursorRead(ensure.Consumer).Do(ctx, tx)
		require.NoError(t, err)
		require.NotNil(t, res)

		require.NoError(t, tx.Commit(ctx))
	})
}
