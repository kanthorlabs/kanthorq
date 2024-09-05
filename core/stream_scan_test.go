package core

import (
	"context"
	"testing"
	"time"

	"github.com/kanthorlabs/kanthorq/pkg/xfaker"
	"github.com/kanthorlabs/kanthorq/tester"
	"github.com/stretchr/testify/require"
)

func TestStreamScan(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	stream, consumer := Seed(t, ctx, conn)
	_ = SeedEvents(t, ctx, conn, stream, consumer, xfaker.F.IntBetween(101, 999))

	req := &StreamScanReq{
		Stream:      stream,
		Consumer:    consumer,
		Size:        xfaker.F.IntBetween(10, 100),
		WaitingTime: time.Second,
	}
	res, err := Do(ctx, req, conn)
	require.NoError(t, err)

	require.Equal(t, req.Size, len(res.Ids))
}

func TestStreamScan_Validate(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	stream, consumer := Seed(t, ctx, conn)
	_ = SeedEvents(t, ctx, conn, stream, consumer, xfaker.F.IntBetween(101, 999))

	req := &StreamScanReq{
		Stream:   stream,
		Consumer: consumer,
		Size:     xfaker.F.IntBetween(10, 100),
	}
	_, err = Do(ctx, req, conn)
	require.Error(t, err)
}

func TestStreamScan_Excludes(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	stream, consumer := Seed(t, ctx, conn)

	// insert excluded events
	events := tester.FakeEvents(xfaker.SubjectWihtPattern(consumer.SubjectExcludes[0]), xfaker.F.IntBetween(101, 999))
	_, err = Do(ctx, &StreamPutEventsReq{Stream: stream, Events: events}, conn)
	require.NoError(t, err)

	req := &StreamScanReq{
		Stream:      stream,
		Consumer:    consumer,
		Size:        xfaker.F.IntBetween(10, 100),
		WaitingTime: time.Second,
	}
	res, err := Do(ctx, req, conn)
	require.NoError(t, err)

	require.Equal(t, 0, len(res.Ids))
}

func TestStreamScan_ExactReqSize(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	stream, consumer := Seed(t, ctx, conn)
	events := SeedEvents(t, ctx, conn, stream, consumer, xfaker.F.IntBetween(10, 100))

	req := &StreamScanReq{
		Stream:      stream,
		Consumer:    consumer,
		Size:        len(events),
		WaitingTime: time.Second,
	}
	res, err := Do(ctx, req, conn)
	require.NoError(t, err)

	require.Equal(t, req.Size, len(res.Ids))
}

func TestStreamScan_LessThanReqSize(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	stream, consumer := Seed(t, ctx, conn)
	events := SeedEvents(t, ctx, conn, stream, consumer, xfaker.F.IntBetween(11, 99))

	req := &StreamScanReq{
		Stream:      stream,
		Consumer:    consumer,
		Size:        len(events) + 1,
		WaitingTime: time.Second,
	}
	res, err := Do(ctx, req, conn)
	require.NoError(t, err)

	require.Equal(t, len(events), len(res.Ids))
}

func TestStreamScan_MixEventSubjects(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	count := xfaker.F.IntBetween(11, 99)

	stream, consumer := Seed(t, ctx, conn)
	// first, seed event with given subject pattern again
	_ = SeedEvents(t, ctx, conn, stream, consumer, count/2+1)
	// then seed events with other subjects
	events := tester.FakeEvents(xfaker.Subject(), xfaker.F.IntBetween(100, 500))
	_, err = Do(ctx, &StreamPutEventsReq{Stream: stream, Events: events}, conn)
	require.NoError(t, err)
	// finally, seed event with given subject pattern again
	_ = SeedEvents(t, ctx, conn, stream, consumer, count/2+1)

	req := &StreamScanReq{
		Stream:      stream,
		Consumer:    consumer,
		Size:        count,
		WaitingTime: time.Second,
	}
	res, err := Do(ctx, req, conn)
	require.NoError(t, err)

	require.Equal(t, count, len(res.Ids))
}

func TestStreamScan_NoEvent(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	stream, consumer := Seed(t, ctx, conn)

	req := &StreamScanReq{
		Stream:      stream,
		Consumer:    consumer,
		Size:        xfaker.F.IntBetween(10, 100),
		WaitingTime: time.Second,
	}
	res, err := Do(ctx, req, conn)
	require.NoError(t, err)

	require.Equal(t, 0, len(res.Ids))
}
