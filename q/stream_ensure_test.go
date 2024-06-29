package q

import (
	"context"
	"testing"
	"time"

	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/kanthorlabs/kanthorq/utils"
	"github.com/stretchr/testify/require"
)

func TestNewStreamEnsure(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ctx := context.Background()
		conn, err := testify.SetupPostgres(ctx)
		require.NoError(t, err)
		defer conn.Close(ctx)

		tx, err := conn.Begin(ctx)
		require.NoError(t, err)

		s, err := NewStreamEnsure(testify.StreamName(5)).Do(ctx, tx)
		require.NoError(t, err)
		require.NotNil(t, s)
		require.NotNil(t, s.Stream)

		require.NoError(t, tx.Commit(ctx))
	})

	t.Run("error of dead connection", func(t *testing.T) {
		ctx := context.Background()
		conn, err := testify.SetupPostgres(ctx)
		require.NoError(t, err)

		tx, err := conn.Begin(ctx)
		require.NoError(t, err)

		conn.Close(ctx)
		_, err = NewStreamEnsure(testify.StreamName(5)).Do(ctx, tx)
		require.ErrorContains(t, err, "conn closed")
	})

	t.Run("error of stream creation", func(t *testing.T) {
		ctx := context.Background()
		name := testify.StreamName(5)
		lock := utils.AdvisoryLockHash(name)

		// simulate deadlock so we cannot create a stream on time
		commit := testify.SetupPostgresAdvisoryXactLock(ctx, lock)
		defer commit()

		timeoutctx, cancel := context.WithTimeout(ctx, time.Second*3)
		defer cancel()

		conn, err := testify.SetupPostgres(ctx)
		require.NoError(t, err)
		defer conn.Close(ctx)

		tx, err := conn.Begin(ctx)
		require.NoError(t, err)
		_, err = NewStreamEnsure(name).Do(timeoutctx, tx)
		require.ErrorIs(t, err, context.DeadlineExceeded)
	})
}
