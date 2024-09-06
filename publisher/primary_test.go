package publisher

import (
	"context"
	"os"
	"testing"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/xfaker"
	"github.com/kanthorlabs/kanthorq/pkg/xlogger"
	"github.com/kanthorlabs/kanthorq/tester"
	"github.com/stretchr/testify/require"
)

func TestPublisher_Validate(t *testing.T) {
	conn, err := tester.SetupPostgres(context.Background())
	require.NoError(t, err)
	defer conn.Close(context.Background())

	options := &Options{
		Connection: os.Getenv("KANTHORQ_POSTGRES_URI"),
	}
	_, err = New(options, xlogger.NewNoop())
	require.Error(t, err)
}

func TestPublisher_Send(t *testing.T) {
	conn, err := tester.SetupPostgres(context.Background())
	require.NoError(t, err)
	defer conn.Close(context.Background())

	options := &Options{
		Connection: os.Getenv("KANTHORQ_POSTGRES_URI"),
		StreamName: entities.DefaultStreamName,
	}
	pub, err := New(options, xlogger.NewNoop())
	require.NoError(t, err)
	require.NoError(t, pub.Start(context.Background()))
	defer func() {
		require.NoError(t, pub.Stop(context.Background()))
	}()

	events := []*entities.Event{
		entities.NewEvent(xfaker.Subject(), []byte("{\"ping\": true}")),
	}
	require.NoError(t, pub.Send(context.Background(), events))
	// send duplicated event
	require.ErrorContains(t, pub.Send(context.Background(), events), "SQLSTATE 23505")
}

func TestPublisher_Send_Validate(t *testing.T) {
	conn, err := tester.SetupPostgres(context.Background())
	require.NoError(t, err)
	defer conn.Close(context.Background())

	options := &Options{
		Connection: os.Getenv("KANTHORQ_POSTGRES_URI"),
		StreamName: entities.DefaultStreamName,
	}
	pub, err := New(options, xlogger.NewNoop())
	require.NoError(t, err)
	require.NoError(t, pub.Start(context.Background()))
	defer func() {
		require.NoError(t, pub.Stop(context.Background()))
	}()

	// no event at all
	events := make([]*entities.Event, 0)
	require.ErrorContains(t, pub.Send(context.Background(), events), "PUBLISHER.SEND.NO_EVENTS")

	// one event has error
	events = append(events, entities.NewEvent(xfaker.Subject(), nil))
	require.ErrorContains(t, pub.Send(context.Background(), events), "Field validation for")
}

func TestPublisher_SendTx(t *testing.T) {
	conn, err := tester.SetupPostgres(context.Background())
	require.NoError(t, err)
	defer conn.Close(context.Background())

	options := &Options{
		Connection: os.Getenv("KANTHORQ_POSTGRES_URI"),
		StreamName: entities.DefaultStreamName,
	}
	pub, err := New(options, xlogger.NewNoop())
	require.NoError(t, err)
	require.NoError(t, pub.Start(context.Background()))
	defer func() {
		require.NoError(t, pub.Stop(context.Background()))
	}()

	events := []*entities.Event{
		entities.NewEvent(xfaker.Subject(), []byte("{\"ping\": true}")),
	}

	tx, err := conn.Begin(context.Background())
	require.NoError(t, err)
	require.NoError(t, pub.SendTx(context.Background(), events, tx))
	require.NoError(t, tx.Commit(context.Background()))
}
