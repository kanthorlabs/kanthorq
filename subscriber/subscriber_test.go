package subscriber

import (
	"context"
	"os"
	"testing"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/puller"
	"github.com/kanthorlabs/kanthorq/tester"
	"github.com/stretchr/testify/require"
)

func TestSubscriber_Connection(t *testing.T) {
	conn, err := tester.SetupPostgres(context.Background())
	require.NoError(t, err)
	defer conn.Close(context.Background())

	instance, err := New(
		&Options{
			Connection:                os.Getenv("KANTHORQ_POSTGRES_URI"),
			StreamName:                entities.DefaultStreamName,
			ConsumerName:              "internal",
			ConsumerSubjectFilter:     []string{"system.ping"},
			ConsumerAttemptMax:        entities.DefaultConsumerAttemptMax,
			ConsumerVisibilityTimeout: entities.DefaultConsumerVisibilityTimeout,
			Puller: puller.PullerIn{
				Size:        100,
				WaitingTime: 3000,
			},
		},
	)
	require.NoError(t, err)
	require.NotNil(t, instance)

	require.NoError(t, instance.Start(context.Background()))
	require.NotNil(t, instance.(*primary).stream, "stream should not be nil")
	require.NotNil(t, instance.(*primary).consumer, "consumer should not be nil")
	require.Equal(t, instance.(*primary).stream.Name, entities.DefaultStreamName, "should use default stream name")
	require.Equal(t, instance.(*primary).consumer.AttemptMax, entities.DefaultConsumerAttemptMax, "should use default consumer attempt max")

	require.NoError(t, instance.Stop(context.Background()))
	require.Nil(t, instance.(*primary).stream, "stream must be deleted after stop")
	require.Nil(t, instance.(*primary).consumer, "consumer must be deleted after stop")
}
