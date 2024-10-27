package core

import (
	"context"
	"testing"

	"github.com/kanthorlabs/kanthorq/pkg/xfaker"
	"github.com/kanthorlabs/kanthorq/tester"
	"github.com/stretchr/testify/require"
)

func TestStreamPutEvents(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	stream, err := Do(ctx, &StreamRegisterReq{StreamName: xfaker.StreamName()}, conn)
	require.NoError(t, err)

	events := tester.FakeEvents(xfaker.Subject(), xfaker.F.IntBetween(100, 200))

	req := &StreamPutEventsReq{Stream: stream.StreamRegistry, Events: events}
	res, err := Do(ctx, req, conn)
	require.NoError(t, err)

	require.Equal(t, int64(len(events)), res.InsertCount)
}

func TestStreamPutEvents_Validate(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	stream, err := Do(ctx, &StreamRegisterReq{StreamName: xfaker.StreamName()}, conn)
	require.NoError(t, err)

	req := &StreamPutEventsReq{Stream: stream.StreamRegistry}
	_, err = Do(ctx, req, conn)
	require.Error(t, err)
}
