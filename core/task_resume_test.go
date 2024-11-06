package core

import (
	"context"
	"testing"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/xfaker"
	"github.com/kanthorlabs/kanthorq/tester"
	"github.com/stretchr/testify/require"
)

func TestTaskResume(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	stream, consumer := Seed(t, ctx, conn)

	events := SeedEvents(t, ctx, conn, stream, consumer, xfaker.F.IntBetween(100, 200))
	tasks := SeedTasks(t, ctx, conn, consumer, events, entities.StateCancelled)

	noopEvents := SeedEvents(t, ctx, conn, stream, consumer, xfaker.F.IntBetween(100, 200))
	noopTasks := SeedTasks(t, ctx, conn, consumer, noopEvents, entities.StatePending)

	req := &TaskResumeReq{
		Consumer: consumer,
		Tasks:    append(tasks, noopTasks...),
	}
	res, err := Do(ctx, conn, req)
	require.NoError(t, err)

	require.Equal(t, len(tasks), len(res.Updated))
	require.Equal(t, len(noopTasks), len(res.Noop))
}

func TestTaskResume_Validate(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	_, consumer := Seed(t, ctx, conn)

	req := &TaskResumeReq{
		Consumer: consumer,
	}
	_, err = Do(ctx, conn, req)
	require.Error(t, err)
}