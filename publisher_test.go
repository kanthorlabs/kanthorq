package kanthorq

import (
	"context"
	"os"
	"testing"

	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/stretchr/testify/require"
)

func TestPublisher_Connection(t *testing.T) {
	conn, err := testify.SetupPostgres(context.Background())
	require.NoError(t, err)
	defer conn.Close(context.Background())

	instance, err := NewPublisher(os.Getenv("KANTHORQ_POSTGRES_URI"))
	require.NoError(t, err)
	require.NotNil(t, instance)

	require.NoError(t, instance.Start(context.Background()))
	require.NotNil(t, instance.(*publisher).conn)
	require.NotNil(t, instance.(*publisher).stream)
	require.Equal(t, instance.(*publisher).stream.Name, DefaultStreamName)

	require.NoError(t, instance.Stop(context.Background()))
	require.Nil(t, instance.(*publisher).conn)
	require.Nil(t, instance.(*publisher).stream)
}

func TestPublisher_Send(t *testing.T) {
	conn, err := testify.SetupPostgres(context.Background())
	require.NoError(t, err)
	defer conn.Close(context.Background())

	instance, err := NewPublisher(os.Getenv("KANTHORQ_POSTGRES_URI"))
	require.NoError(t, err)
	require.NotNil(t, instance)
	require.NoError(t, instance.Start(context.Background()))
	defer func() {
		require.NoError(t, instance.Stop(context.Background()))
	}()

	event := NewEvent(testify.Topic(1), []byte("{\"ping\": true}"))
	require.NoError(t, instance.Send(context.Background(), event))
}
