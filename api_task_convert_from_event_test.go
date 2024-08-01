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
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	_, c, events := FakeEntities(t, ctx, conn, faker.F.IntBetween(100, 500))
	remaining := faker.F.IntBetween(1, 9)
	size := len(events) - remaining

	first, err := Do(ctx, &TaskConvertFromEventReq{
		Consumer:         c,
		InitialTaskState: StateAvailable,
		Size:             size,
		ScanWindow:       30000,
		ScanRoundMax:     3,
		ScanRoundDelay:   1000,
	}, conn)
	require.NoError(t, err)

	require.Equal(t, size, len(first.Tasks))
	require.Equal(t, size, len(first.EventIds))

	seconds, err := Do(ctx, &TaskConvertFromEventReq{
		Consumer:         c,
		InitialTaskState: StateAvailable,
		Size:             size,
		ScanWindow:       30000,
		ScanRoundMax:     3,
		ScanRoundDelay:   1000,
	}, conn)
	require.NoError(t, err)

	require.Equal(t, remaining, len(seconds.Tasks))
	require.Equal(t, remaining, len(seconds.EventIds))
}
