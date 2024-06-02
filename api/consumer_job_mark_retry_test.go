package api

import (
	"context"
	"errors"
	"testing"

	"github.com/kanthorlabs/common/clock"
	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/stretchr/testify/require"
)

func TestConsumerJobMarkRetry(t *testing.T) {
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

		events := testify.GenStreamEvents(ctx, testify.StreamName(5), 10)
		var maps = make(map[string]error, 0)
		for _, event := range events {
			maps[event.EventId] = errors.New(event.EventId)
		}

		r, err := ConsumerJobMarkRetry(c.Consumer, maps).Do(ctx, tx, clock.New())
		require.NoError(t, err)
		require.Equal(t, len(events), len(r.Status))

		require.NoError(t, tx.Commit(ctx))
	})
}
