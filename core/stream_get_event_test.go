package core

import (
	"context"
	"testing"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/faker"
	"github.com/kanthorlabs/kanthorq/tester"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestEventGet(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	stream, consumer := Seed(t, ctx, conn)
	events := SeedEvents(t, ctx, conn, stream, consumer, faker.F.IntBetween(100, 500))

	req := &StreamGetEventReq{
		Stream:   stream,
		EventIds: lo.Map(events, func(e *entities.Event, _ int) string { return e.Id }),
	}
	res, err := Do(ctx, req, conn)
	require.NoError(t, err)

	require.Equal(t, len(events), len(res.Events))
}
