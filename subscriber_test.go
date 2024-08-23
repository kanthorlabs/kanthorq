package kanthorq

import (
	"context"
	"os"
	"testing"

	"github.com/kanthorlabs/kanthorq/tester"
	"github.com/stretchr/testify/require"
)

func TestSubscriber_Connection(t *testing.T) {
	conn, err := tester.SetupPostgres(context.Background())
	require.NoError(t, err)
	defer conn.Close(context.Background())

	instance, err := NewSubscriber(
		os.Getenv("KANTHORQ_POSTGRES_URI"),
		&SubscriberOptions{
			StreamName:            DefaultStreamName,
			ConsumerName:          "internal",
			ConsumerSubjectFilter: []string{"system.ping"},
			ConsumerAttemptMax:    DefaultConsumerAttemptMax,
		},
	)
	require.NoError(t, err)
	require.NotNil(t, instance)

	require.NoError(t, instance.Start(context.Background()))
	require.NotNil(t, instance.(*subscriber).stream, "stream should not be nil")
	require.NotNil(t, instance.(*subscriber).consumer, "consumer should not be nil")
	require.Equal(t, instance.(*subscriber).stream.Name, DefaultStreamName, "should use default stream name")
	require.Equal(t, instance.(*subscriber).consumer.AttemptMax, DefaultConsumerAttemptMax, "should use default consumer attempt max")

	require.NoError(t, instance.Stop(context.Background()))
	require.Nil(t, instance.(*subscriber).stream, "stream must be deleted after stop")
	require.Nil(t, instance.(*subscriber).consumer, "consumer must be deleted after stop")
}
