package api

import (
	"context"
	"testing"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/stretchr/testify/require"
)

func TestConsumerEnsure(t *testing.T) {
	t.Run("happy case", func(t *testing.T) {
		ctx := context.Background()

		pool, err := testify.SetupPostgres(ctx)
		require.NoError(t, err)

		tx, err := pool.Begin(ctx)
		require.NoError(t, err)

		ensure, err := ConsumerEnsure(
			&entities.Stream{Name: testify.Fake.RandomStringWithLength(32)},
			testify.Fake.RandomStringWithLength(32),
			testify.Fake.RandomStringWithLength(32),
		).Do(ctx, tx)
		require.NoError(t, err)
		require.NotNil(t, ensure)
		require.NotNil(t, ensure.Consumer)

		require.NoError(t, tx.Commit(ctx))
	})
}
