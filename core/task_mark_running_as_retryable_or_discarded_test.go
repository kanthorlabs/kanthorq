package core

import (
	"context"
	"runtime/debug"
	"testing"
	"time"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/xfaker"
	"github.com/kanthorlabs/kanthorq/tester"
	"github.com/stretchr/testify/require"
)

func TestTaskMarkRunningAsRetryableOrDiscarded_ToRetryable(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	stream, consumer := Seed(t, ctx, conn)

	events := SeedEvents(t, ctx, conn, stream, consumer, xfaker.F.IntBetween(100, 200))
	tasks := SeedTasks(t, ctx, conn, consumer, events, entities.StateRunning)

	noopEvents := SeedEvents(t, ctx, conn, stream, consumer, xfaker.F.IntBetween(100, 200))
	noopTasks := SeedTasks(t, ctx, conn, consumer, noopEvents, entities.StateCancelled)

	req := &TaskMarkRunningAsRetryableOrDiscardedReq{
		Consumer: consumer,
		Tasks:    append(tasks, noopTasks...),
		Error: entities.AttemptedError{
			At:    time.Now().UnixMilli(),
			Error: "error",
			Stack: string(debug.Stack()),
		},
	}
	res, err := Do(ctx, conn, req)
	require.NoError(t, err)

	require.Equal(t, len(tasks), len(res.Updated))
	require.Equal(t, len(tasks), len(res.States))
	require.Equal(t, len(noopTasks), len(res.Noop))

	for _, state := range res.States {
		require.Equal(t, entities.StateRetryable, state)
	}
}

func TestTaskMarkRunningAsRetryableOrDiscarded_ToDiscarded(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	stream, consumer := Seed(t, ctx, conn)

	events := SeedEvents(t, ctx, conn, stream, consumer, xfaker.F.IntBetween(100, 200))
	tasks := tester.FakeTasks(events, entities.StateRunning)
	for i := range tasks {
		// simulate that we have reached the max attempts
		tasks[i].AttemptCount = consumer.AttemptMax + 1
	}
	_, err = Do(ctx, conn, &ConsumerPutTasksReq{Consumer: consumer, Tasks: tasks})
	require.NoError(t, err)

	noopEvents := SeedEvents(t, ctx, conn, stream, consumer, xfaker.F.IntBetween(100, 200))
	noopTasks := SeedTasks(t, ctx, conn, consumer, noopEvents, entities.StateCancelled)

	req := &TaskMarkRunningAsRetryableOrDiscardedReq{
		Consumer: consumer,
		Tasks:    append(tasks, noopTasks...),
		Error: entities.AttemptedError{
			At:    time.Now().UnixMilli(),
			Error: "error",
			Stack: string(debug.Stack()),
		},
	}
	res, err := Do(ctx, conn, req)
	require.NoError(t, err)

	require.Equal(t, len(tasks), len(res.Updated))
	require.Equal(t, len(tasks), len(res.States))
	require.Equal(t, len(noopTasks), len(res.Noop))

	for _, state := range res.States {
		require.Equal(t, entities.StateDiscarded, state)
	}
}

func TestTaskMarkRunningAsRetryableOrDiscarded_Validate(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	stream, consumer := Seed(t, ctx, conn)

	events := SeedEvents(t, ctx, conn, stream, consumer, xfaker.F.IntBetween(100, 200))
	tasks := SeedTasks(t, ctx, conn, consumer, events, entities.StateRunning)

	noopEvents := SeedEvents(t, ctx, conn, stream, consumer, xfaker.F.IntBetween(100, 200))
	noopTasks := SeedTasks(t, ctx, conn, consumer, noopEvents, entities.StateCancelled)

	req := &TaskMarkRunningAsRetryableOrDiscardedReq{
		Consumer: consumer,
		Tasks:    append(tasks, noopTasks...),
		Error: entities.AttemptedError{
			// missing At
			Error: "error",
			Stack: string(debug.Stack()),
		},
	}
	_, err = Do(ctx, conn, req)
	require.Error(t, err)
}
