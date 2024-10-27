package core

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/pkg/xfaker"
	"github.com/kanthorlabs/kanthorq/pkg/xid"
	"github.com/kanthorlabs/kanthorq/tester"
	"github.com/stretchr/testify/require"
)

func TestConsumerUnlock(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	_, consumer := Seed(t, ctx, conn)

	req := &ConsumerUnlockReq{
		Name:   consumer.Name,
		Cursor: xid.New("test"),
	}
	res, err := Do(ctx, conn, req)
	require.NoError(t, err)

	require.Equal(t, req.Cursor, res.Consumer.Cursor)
}

func TestConsumerUnlock_Validate(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	_, consumer := Seed(t, ctx, conn)

	req := &ConsumerUnlockReq{
		Name: consumer.Name,
	}
	_, err = Do(ctx, conn, req)
	require.Error(t, err)
}

func TestConsumerUnlock_Failure(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	_, _ = Seed(t, ctx, conn)

	req := &ConsumerUnlockReq{
		Name:   xfaker.ConsumerName(),
		Cursor: xid.New("test"),
	}
	_, err = Do(ctx, conn, req)
	require.ErrorIs(t, err, pgx.ErrNoRows)
}
