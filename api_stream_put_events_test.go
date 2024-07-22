package kanthorq

import (
	"context"
	"testing"

	"github.com/kanthorlabs/kanthorq/pkg/faker"
	"github.com/kanthorlabs/kanthorq/tester"
	"github.com/stretchr/testify/require"
)

func TestStreamPut(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	require.NoError(t, err)

	stream, err := StreamRegister(ctx, &StreamRegisterReq{StreamName: faker.StreamName()}, conn)
	require.NoError(t, err)

	count := faker.F.Int64Between(100, 500)
	events := make([]*Event, count)
	for i := 0; i < int(count); i++ {
		events[i] = NewEvent(faker.Topic(), faker.DataOf16Kb())
	}

	req := &StreamPutEventsReq{Stream: stream.StreamRegistry, Events: events}
	res, err := StreamPutEvents(ctx, req, conn)
	require.NoError(t, err)

	require.Equal(t, count, res.InsertCount)
}
