package core

import (
	"context"
	"testing"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/xfaker"
	"github.com/kanthorlabs/kanthorq/tester"
	"github.com/stretchr/testify/require"
)

func TestTaskStateTransition(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	stream, consumer := Seed(t, ctx, conn)

	events := SeedEvents(t, ctx, conn, stream, consumer, xfaker.F.IntBetween(100, 500))
	tasks := SeedTasks(t, ctx, conn, consumer, events, entities.StateRetryable)

	noopEvents := SeedEvents(t, ctx, conn, stream, consumer, xfaker.F.IntBetween(100, 500))
	noopTasks := SeedTasks(t, ctx, conn, consumer, noopEvents, entities.StateCancelled)

	req := &TaskStateTransitionReq{
		Consumer:  consumer,
		FromState: entities.StateRetryable,
		ToState:   entities.StateRunning,
		Size:      len(tasks) + len(noopTasks),
	}
	res, err := Do(ctx, req, conn)
	require.NoError(t, err)

	require.Equal(t, len(tasks), len(res.Tasks))
}