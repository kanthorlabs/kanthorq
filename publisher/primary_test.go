package publisher

import (
	"context"
	"os"
	"testing"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/xfaker"
	"github.com/kanthorlabs/kanthorq/tester"
	"github.com/stretchr/testify/require"
)

func TestPublisher_Validate(t *testing.T) {
	conn, err := tester.SetupPostgres(context.Background())
	require.NoError(t, err)
	defer conn.Close(context.Background())

	_, err = New(&Options{
		Connection: os.Getenv("KANTHORQ_POSTGRES_URI"),
	})
	require.Error(t, err)
}

func TestPublisher_Send(t *testing.T) {
	conn, err := tester.SetupPostgres(context.Background())
	require.NoError(t, err)
	defer conn.Close(context.Background())

	instance, err := New(&Options{
		Connection: os.Getenv("KANTHORQ_POSTGRES_URI"),
		StreamName: entities.DefaultStreamName,
	})
	require.NoError(t, err)
	require.NoError(t, instance.Start(context.Background()))
	defer func() {
		require.NoError(t, instance.Stop(context.Background()))
	}()

	events := []*entities.Event{
		entities.NewEvent(xfaker.Subject(), []byte("{\"ping\": true}")),
	}
	require.NoError(t, instance.Send(context.Background(), events))
	// send duplicated event
	require.ErrorContains(t, instance.Send(context.Background(), events), "SQLSTATE 23505")
}

func TestPublisher_Send_Validate(t *testing.T) {
	conn, err := tester.SetupPostgres(context.Background())
	require.NoError(t, err)
	defer conn.Close(context.Background())

	instance, err := New(&Options{
		Connection: os.Getenv("KANTHORQ_POSTGRES_URI"),
		StreamName: entities.DefaultStreamName,
	})
	require.NoError(t, err)
	require.NoError(t, instance.Start(context.Background()))
	defer func() {
		require.NoError(t, instance.Stop(context.Background()))
	}()

	events := make([]*entities.Event, 0)
	require.ErrorContains(t, instance.Send(context.Background(), events), "no events provided")
	events = append(events, entities.NewEvent(xfaker.Subject(), nil))
	require.ErrorContains(t, instance.Send(context.Background(), events), "Field validation for")
}

func TestPublisher_SendTx(t *testing.T) {
	conn, err := tester.SetupPostgres(context.Background())
	require.NoError(t, err)
	defer conn.Close(context.Background())

	instance, err := New(&Options{
		Connection: os.Getenv("KANTHORQ_POSTGRES_URI"),
		StreamName: entities.DefaultStreamName,
	})
	require.NoError(t, err)
	require.NoError(t, instance.Start(context.Background()))
	defer func() {
		require.NoError(t, instance.Stop(context.Background()))
	}()

	events := []*entities.Event{
		entities.NewEvent(xfaker.Subject(), []byte("{\"ping\": true}")),
	}

	tx, err := conn.Begin(context.Background())
	require.NoError(t, err)
	require.NoError(t, instance.SendTx(context.Background(), events, tx))
	require.NoError(t, tx.Commit(context.Background()))
}
