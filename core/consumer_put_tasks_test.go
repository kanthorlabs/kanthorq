package core

import (
	"context"
	"testing"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/xfaker"
	"github.com/kanthorlabs/kanthorq/tester"
	"github.com/stretchr/testify/require"
)

func TestConsumerPutTasks(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	stream, consumer := Seed(t, ctx, conn)
	events := SeedEvents(t, ctx, conn, stream, consumer, xfaker.F.IntBetween(100, 200))

	tasks := tester.FakeTasks(events, entities.StateAvailable)
	req := &ConsumerPutTasksReq{
		Consumer: consumer,
		Tasks:    tasks,
	}
	res, err := Do(ctx, conn, req)
	require.NoError(t, err)
	require.Equal(t, int64(len(tasks)), res.InsertCount)
}

func TestConsumerPutTasks_Validate(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	_, consumer := Seed(t, ctx, conn)
	req := &ConsumerPutTasksReq{
		Consumer: consumer,
	}
	_, err = Do(ctx, conn, req)
	require.Error(t, err)
}
