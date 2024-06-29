package q

import (
	"context"
	"testing"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/stretchr/testify/require"
)

func TestNewConsumerCursorRead(t *testing.T) {
	t.Run("happy case", func(t *testing.T) {
		ctx := context.Background()

		pool, err := testify.SetupPostgres(ctx)
		require.NoError(t, err)

		// seed an consumer
		tx, err := pool.Begin(ctx)
		require.NoError(t, err)

		ensure, err := NewConsumerEnsure(
			&entities.Stream{Name: testify.StreamName(5)},
			testify.ConsumerName(5),
			testify.Topic(5),
		).Do(ctx, tx)
		require.NoError(t, err)
		require.NotNil(t, ensure)
		require.NotNil(t, ensure.Consumer)

		// need to commit to make sure the consumer exists in next transaction
		require.NoError(t, tx.Commit(ctx))

		tx, err = pool.Begin(ctx)
		require.NoError(t, err)

		res, err := NewConsumerCursorRead(ensure.Consumer).Do(ctx, tx)
		require.NoError(t, err)
		require.NotNil(t, res)

		require.NoError(t, tx.Commit(ctx))
	})
}
