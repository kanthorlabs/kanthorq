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

	events := FakeEvents(faker.Topic(), 100, 500)

	req := &StreamPutEventsReq{Stream: stream.StreamRegistry, Events: events}
	res, err := StreamPutEvents(ctx, req, conn)
	require.NoError(t, err)

	require.Equal(t, int64(len(events)), res.InsertCount)
}
