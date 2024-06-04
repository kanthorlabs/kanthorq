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

	// publish first
	pub, err := publisher.New(ctx, &publisher.Config{
		ConnectionUri: os.Getenv("TEST_DATABASE_URI"),
		StreamName:    streamName,
	})
	require.NoError(t, err)

	events := testify.GenStreamEvents(ctx, topic, 1000)
	require.NoError(t, pub.Send(ctx, events))

	// then subscribe
	sub, err := New(ctx, &Config{
		ConnectionUri: os.Getenv("TEST_DATABASE_URI"),
		StreamName:    streamName,
		ConsumerName:  consumerName,
		Topic:         topic,
	})

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
	case err, ok := <-sub.Error():
		if !ok {
			panic("oops, something went wrong")
		}
		cancel()
		require.NoError(t, err)
	case reports, ok := <-sub.Report():
		if !ok {
			panic("oops, something went wrong")
		}
		cancel()
		require.Greater(t, len(reports), 0)
		require.Less(t, len(reports), DefaultSize)
	}
}
