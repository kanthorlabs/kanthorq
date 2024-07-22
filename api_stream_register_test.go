package kanthorq

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/pkg/faker"
	"github.com/kanthorlabs/kanthorq/tester"
	"github.com/stretchr/testify/require"
)

func TestStreamRegister(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	require.NoError(t, err)

	req := &StreamRegisterReq{
		StreamName: faker.StreamName(),
	}
	res, err := StreamRegister(ctx, req, conn)
	require.NoError(t, err)

	require.NotNil(t, res)
	require.Equal(t, req.StreamName, res.StreamRegistry.Name)
}

func TestStreamRegister_Parallel(t *testing.T) {
	ctx := context.Background()

	count := faker.F.IntBetween(20, 30)
	// setup all connections so we don't waste time on it during the test
	conns := make([]*pgx.Conn, count)
	for i := 0; i < count; i++ {
		var err error
		conns[i], err = tester.SetupPostgres(ctx)
		require.NoError(t, err)
	}

	// will try to register same stream
	req := &StreamRegisterReq{
		StreamName: faker.StreamName(),
	}

	for i := 0; i < count; i++ {
		conn := conns[i]

		t.Run(fmt.Sprintf("parallel #%d", i), func(subt *testing.T) {
			subt.Parallel()
			res, err := StreamRegister(ctx, req, conn)
			require.NoError(subt, err)

			require.NotNil(subt, res)
			require.Equal(subt, req.StreamName, res.StreamRegistry.Name)
		})
	}
}
