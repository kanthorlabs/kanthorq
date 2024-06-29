package q

import (
	"context"
	"testing"
	"time"

	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/kanthorlabs/kanthorq/utils"
	"github.com/stretchr/testify/require"
)

func TestNewConsumerJobCreate(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ctx := context.Background()

		conn, err := testify.SetupPostgres(ctx)
		require.NoError(t, err)
		defer conn.Close(ctx)

		tx, err := conn.Begin(ctx)
		require.NoError(t, err)

		err = NewConsumerJobCreate(testify.ConsumerName(5)).Do(ctx, tx)
		require.NoError(t, err)

		require.NoError(t, tx.Commit(ctx))
	})

	t.Run("should able to handle concurrency creation", func(t *testing.T) {
		ctx := context.Background()
		name := testify.ConsumerName(5)
		lock := utils.AdvisoryLockHash(name)

		lconn, err := testify.SetupPostgres(ctx)
		require.NoError(t, err)
		defer lconn.Close(ctx)

		ltx, err := lconn.Begin(ctx)
		require.NoError(t, err)
		_, err = ltx.Exec(ctx, "SELECT pg_advisory_xact_lock($1)", lock)
		require.NoError(t, err)

		go func() {
			conn, err := testify.SetupPostgres(ctx)
			require.NoError(t, err)
			defer conn.Close(ctx)

			tx, err := conn.Begin(ctx)
			require.NoError(t, err)
			require.NoError(t, NewConsumerJobCreate(name).Do(ctx, tx))
			require.NoError(t, tx.Commit(ctx))
		}()

		// wait a bit to ensure the go routine has started
		time.Sleep(time.Second * 3)

		require.NoError(t, ltx.Commit(ctx))
	})
}
