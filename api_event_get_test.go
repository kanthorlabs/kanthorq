package kanthorq

import (
	"context"
	"testing"

	"github.com/kanthorlabs/kanthorq/pkg/faker"
	"github.com/kanthorlabs/kanthorq/tester"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestEventGet(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	require.NoError(t, err)

	s, _, events := FakeEntities(t, ctx, conn, faker.F.IntBetween(100, 500))

	req := &EventGetReq{
		Stream:   s,
		EventIds: lo.Map(events, func(e *Event, _ int) string { return e.Id }),
	}
	res, err := Do(ctx, req, conn)
	require.NoError(t, err)

	require.Equal(t, len(events), len(res.Events))
}
