package kanthorq

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/stretchr/testify/require"
)

func TestSubscriber_Connection(t *testing.T) {
	conn, err := testify.SetupPostgres(context.Background())
	require.NoError(t, err)
	defer conn.Close(context.Background())

	instance, err := NewSubscriber(
		os.Getenv("KANTHORQ_POSTGRES_URI"),
		&SubscriberOptions{
			StreamName:         DefaultStreamName,
			ConsumerName:       "internal",
			ConsumerTopic:      "system.ping",
			ConsumerAttemptMax: DefaultConsumerAttemptMax,
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

func handle(ctx context.Context, fn func()) error {
	delay := 1000
	for {
		hctx, cancel := context.WithTimeout(ctx, time.Millisecond*12000)
		defer cancel()

		select {
		case <-hctx.Done():
			return hctx.Err()
		default:
			fn()

			fmt.Println("---->", delay)
			delay += 1000
			timer := time.NewTimer(time.Millisecond * time.Duration(delay))

			select {
			case <-timer.C:
				// do nothing, just wait
			case <-hctx.Done():
				timer.Stop()
				// if context got canceled, should stop both the loop and the delay
				return hctx.Err()
			}
		}
	}
}

func TestIdea(t *testing.T) {
	ctx := context.Background()
	err := handle(ctx, func() { time.Sleep(time.Millisecond * 7000) })
	require.NoError(t, err)
}
