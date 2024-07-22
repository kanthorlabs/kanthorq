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

func TestConsumerRegister(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	require.NoError(t, err)

	req := &ConsumerRegisterReq{
		ConsumerName: faker.ConsumerName(),
	}
	res, err := ConsumerRegister(ctx, req, conn)
	require.NoError(t, err)

	require.NotNil(t, res)
	require.Equal(t, req.ConsumerName, res.ConsumerRegistry.Name)
}

func TestConsumerRegister_Parallel(t *testing.T) {
	ctx := context.Background()

	count := faker.F.IntBetween(20, 30)
	// setup all connections so we don't waste time on it during the test
	conns := make([]*pgx.Conn, count)
	for i := 0; i < count; i++ {
		var err error
		conns[i], err = tester.SetupPostgres(ctx)
		require.NoError(t, err)
	}

	// will try to register same Consumer
	req := &ConsumerRegisterReq{
		ConsumerName: faker.ConsumerName(),
	}

	for i := 0; i < count; i++ {
		conn := conns[i]

		t.Run(fmt.Sprintf("parallel #%d", i), func(subt *testing.T) {
			subt.Parallel()
			res, err := ConsumerRegister(ctx, req, conn)
			require.NoError(subt, err)

			require.NotNil(subt, res)
			require.Equal(subt, req.ConsumerName, res.ConsumerRegistry.Name)
		})
	}
}
