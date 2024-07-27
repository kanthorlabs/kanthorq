package kanthorq

import (
	"context"
	"testing"

	"github.com/kanthorlabs/kanthorq/pkg/faker"
	"github.com/kanthorlabs/kanthorq/tester"
	"github.com/stretchr/testify/require"
)

func TestTaskConvertFromEvent(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	require.NoError(t, err)

	_, c, events := FakeEntities(t, ctx, conn, faker.F.IntBetween(100, 500))
	size := len(events) - 9

	req := &TaskConvertFromEventReq{
		Consumer:         c,
		Size:             size,
		InitialTaskState: StateAvailable,
	}
	res, err := Do(ctx, req, conn)
	require.NoError(t, err)

	require.Equal(t, size, len(res.Tasks))
	require.Equal(t, size, len(res.EventIds))
}
