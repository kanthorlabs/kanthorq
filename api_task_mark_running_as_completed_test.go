package kanthorq

import (
	"context"
	"testing"

	"github.com/kanthorlabs/kanthorq/pkg/faker"
	"github.com/kanthorlabs/kanthorq/tester"
	"github.com/stretchr/testify/require"
)

func TestTaskMarkRunningAsCompleted(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	stream, consumer := Seed(t, ctx, conn)

	events := SeedEvents(t, ctx, conn, stream, consumer, faker.F.IntBetween(100, 500))
	tasks := SeedTasks(t, ctx, conn, consumer, events, StateRunning)

	noopEvents := SeedEvents(t, ctx, conn, stream, consumer, faker.F.IntBetween(100, 500))
	noopTasks := SeedTasks(t, ctx, conn, consumer, noopEvents, StateCancelled)

	req := &TaskMarkRunningAsCompletedReq{
		Consumer: consumer,
		Tasks:    append(tasks, noopTasks...),
	}
	res, err := Do(ctx, req, conn)
	require.NoError(t, err)

	require.Equal(t, len(tasks), len(res.Updated))
	require.Equal(t, len(noopTasks), len(res.Noop))
}
