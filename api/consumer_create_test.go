package api

import (
	"context"
	"testing"

	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/stretchr/testify/require"
)

func TestConsumerCreate(t *testing.T) {
	t.Run("happy case", func(t *testing.T) {
		ctx := context.Background()

		pool, err := testify.SetupPostgres(ctx)
		require.NoError(t, err)

		tx, err := pool.Begin(ctx)
		require.NoError(t, err)

		err = ConsumerCreate(testify.Fake.RandomStringWithLength(32)).Do(ctx, tx)
		require.NoError(t, err)

		require.NoError(t, tx.Commit(ctx))
	})
}