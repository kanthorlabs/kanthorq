package core

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/pkg/xfaker"
	"github.com/kanthorlabs/kanthorq/tester"
	"github.com/stretchr/testify/require"
)

func TestConsumerLock(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	_, consumer := Seed(t, ctx, conn)
	req := &ConsumerLockReq{
		Name: consumer.Name,
	}

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	require.NoError(t, err)
	_, err = req.Do(ctx, tx)
	require.NoError(t, err)

	var count = xfaker.F.IntBetween(11, 21)
	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			conn, err := pgx.Connect(ctx, os.Getenv("KANTHORQ_POSTGRES_URI"))
			require.NoError(t, err)
			defer func() {
				require.NoError(t, conn.Close(ctx))
			}()

			_, err = Do(ctx, req, conn)
			require.ErrorIs(t, err, pgx.ErrNoRows)
		}()
	}
	wg.Wait()

	require.NoError(t, tx.Commit(ctx))
}

func TestConsumerLock_Failure(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	_, _ = Seed(t, ctx, conn)

	req := &ConsumerLockReq{
		Name: xfaker.ConsumerName(),
	}
	_, err = Do(ctx, req, conn)
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestConsumerLock_Validate(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	_, _ = Seed(t, ctx, conn)

	req := &ConsumerLockReq{}
	_, err = Do(ctx, req, conn)
	require.Error(t, err)
}
