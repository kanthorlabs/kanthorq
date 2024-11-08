package core

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/pkg/xfaker"
	"github.com/kanthorlabs/kanthorq/tester"
	"github.com/stretchr/testify/require"
)

func TestStreamRegister(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	req := &StreamRegisterReq{
		StreamName: xfaker.StreamName(),
	}
	res, err := Do(ctx, conn, req)
	require.NoError(t, err)

	require.NotNil(t, res)
	require.Equal(t, req.StreamName, res.StreamRegistry.Name)
}

func TestStreamRegister_Validate(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	req := &StreamRegisterReq{}
	_, err = Do(ctx, conn, req)
	require.Error(t, err)
}

func TestStreamRegister_Parallel(t *testing.T) {
	ctx := context.Background()

	count := xfaker.F.IntBetween(20, 30)
	// setup all connections so we don't waste time on it during the test
	conns := make([]*pgx.Conn, count)
	for i := 0; i < count; i++ {
		var err error
		conns[i], err = tester.SetupPostgres(ctx)
		require.NoError(t, err)
	}

	// will try to register same stream
	req := &StreamRegisterReq{
		StreamName: xfaker.StreamName(),
	}

	for i := 0; i < count; i++ {
		conn := conns[i]

		t.Run(fmt.Sprintf("parallel #%d", i), func(subt *testing.T) {
			subt.Parallel()
			res, err := Do(ctx, conn, req)
			require.NoError(subt, err)

			require.NotNil(subt, res)
			require.Equal(subt, req.StreamName, res.StreamRegistry.Name)
		})
	}
}

func TestStreamRegister_Req(t *testing.T) {
	_, err := (&StreamRegisterReq{}).Do(context.Background(), nil)
	require.Error(t, err)

	validationErrors, ok := err.(validator.ValidationErrors)
	require.True(t, ok)

	require.GreaterOrEqual(t, len(validationErrors), 1)
}
