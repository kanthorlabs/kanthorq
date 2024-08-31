package core

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5"
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

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			conn, err := pgx.Connect(ctx, os.Getenv("KANTHORQ_POSTGRES_URI"))
			require.NoError(t, err)
			defer func() {
				require.NoError(t, conn.Close(ctx))
			}()

			res, err := Do(ctx, req, conn)
			if err != nil {
				require.ErrorIs(t, err, pgx.ErrNoRows)
			}
			if res != nil {
				require.Equal(t, res.Consumer.Name, consumer.Name)
			}
		}()
	}
	wg.Wait()
}
