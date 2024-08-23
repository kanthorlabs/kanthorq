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

func TestPublisher_Connection(t *testing.T) {
	conn, err := tester.SetupPostgres(context.Background())
	require.NoError(t, err)
	defer conn.Close(context.Background())

	instance, err := New(
		&Options{
			Connection: os.Getenv("KANTHORQ_POSTGRES_URI"),
			StreamName: entities.DefaultStreamName,
		},
	)
	require.NoError(t, err)
	require.NotNil(t, instance)

	require.NoError(t, instance.Start(context.Background()))
	require.NotNil(t, instance.(*primary).stream, "stream should not be nil")
	require.Equal(t, instance.(*primary).stream.Name, entities.DefaultStreamName, "should use default stream name")

	require.NoError(t, instance.Stop(context.Background()))
	require.Nil(t, instance.(*primary).stream, "stream must be deleted after stop")
}

func TestPublisher_Send(t *testing.T) {
	conn, err := tester.SetupPostgres(context.Background())
	require.NoError(t, err)
	defer conn.Close(context.Background())

	instance, err := New(
		&Options{
			Connection: os.Getenv("KANTHORQ_POSTGRES_URI"),
			StreamName: entities.DefaultStreamName,
		},
	)
	require.NoError(t, err)
	require.NotNil(t, instance)
	require.NoError(t, instance.Start(context.Background()))
	defer func() {
		require.NoError(t, instance.Stop(context.Background()))
	}()

	event := entities.NewEvent(xfaker.Subject(), []byte("{\"ping\": true}"))
	require.NoError(t, instance.Send(context.Background(), event))
}
