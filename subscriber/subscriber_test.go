package subscriber

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/publisher"
	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/stretchr/testify/require"
)

func TestSubscriber(t *testing.T) {
	ctx := context.Background()

	var streamName = testify.StreamName(5)
	var consumerName = testify.ConsumerName(5)
	var topic = testify.Topic(5)

	pub := publisher.New(&publisher.Config{
		ConnectionUri: os.Getenv("KANTHORQ_POSTGRES_URI"),
		StreamName:    streamName,
	})
	require.NoError(t, pub.Start(ctx))
	defer func() {
		require.NoError(t, pub.Stop(ctx))
	}()

	events := testify.GenStreamEvents(topic, 1000)
	require.NoError(t, pub.Send(ctx, events))

	// then subscribe
	sub := New(&Config{
		ConnectionUri: os.Getenv("KANTHORQ_POSTGRES_URI"),
		StreamName:    streamName,
		ConsumerName:  consumerName,
		Topic:         topic,
	})
	require.NoError(t, sub.Start(ctx))
	defer func() {
		require.NoError(t, sub.Stop(ctx))
	}()

	cancelling, cancel := context.WithCancel(context.Background())
	// receiving events
	go sub.Consume(cancelling, func(ctx context.Context, events map[string]*entities.StreamEvent) map[string]error {
		var reports = make(map[string]error, len(events))

		var i = 1
		for _, event := range events {
			if i%2 == 0 {
				continue
			}

			reports[event.EventId] = errors.New(testify.Fake.Emoji().EmojiCode())
			i++
		}
		return reports
	})

	select {
	case err := <-sub.Error():
		cancel()
		require.NoError(t, err)
	case failures := <-sub.Failurec():
		cancel()
		require.Greater(t, len(failures), 0)
		require.Less(t, len(failures), DefaultSize)
	}
}
