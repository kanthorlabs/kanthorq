package core

import (
	"context"
	"testing"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/xfaker"
	"github.com/kanthorlabs/kanthorq/tester"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestTaskConvert(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	stream, consumer := Seed(t, ctx, conn)

	events := SeedEvents(t, ctx, conn, stream, consumer, xfaker.F.IntBetween(100, 500))

	req := &TaskConvertReq{
		Consumer:     consumer,
		EventIds:     lo.Map(events, func(e *entities.Event, _ int) string { return e.Id }),
		InitialState: entities.StateRunning,
	}
	res, err := Do(ctx, req, conn)
	require.NoError(t, err)

	require.Equal(t, len(events), len(res.Tasks))
	require.Equal(t, len(events), len(res.EventIds))

	// duplicated
	dupres, err := Do(ctx, req, conn)
	require.NoError(t, err)
	require.Equal(t, 0, len(dupres.Tasks))
}

func TestTaskConvert_Validate(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	_, consumer := Seed(t, ctx, conn)
	req := &TaskConvertReq{
		Consumer: consumer,
	}
	_, err = Do(ctx, req, conn)
	require.Error(t, err)
}

func TestTaskConvert_NoEvent(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	_, consumer := Seed(t, ctx, conn)

	events := tester.FakeEvents(xfaker.SubjectWihtPattern(consumer.SubjectIncludes[0]), xfaker.F.IntBetween(100, 500))

	req := &TaskConvertReq{
		Consumer:     consumer,
		EventIds:     lo.Map(events, func(e *entities.Event, _ int) string { return e.Id }),
		InitialState: entities.StateRunning,
	}
	res, err := Do(ctx, req, conn)
	require.NoError(t, err)
	require.Equal(t, 0, len(res.Tasks))
}
