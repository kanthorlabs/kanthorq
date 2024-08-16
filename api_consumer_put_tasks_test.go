package kanthorq

import (
	"context"
	"testing"

	"github.com/kanthorlabs/kanthorq/pkg/faker"
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
	events := SeedEvents(t, ctx, conn, stream, consumer, faker.F.IntBetween(100, 500))

	tasks := FakeTasks(events, StateAvailable)
	req := &ConsumerPutTasksReq{
		Consumer: consumer,
		Tasks:    tasks,
	}
	res, err := Do(ctx, req, conn)
	require.NoError(t, err)
	require.Equal(t, int64(len(tasks)), res.InsertCount)
}
