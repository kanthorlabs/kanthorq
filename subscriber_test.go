package kanthorq

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/kanthorlabs/common/clock"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/stretchr/testify/require"
)

func TestSubscriber(t *testing.T) {
	ctx := context.Background()

	pool, err := Connection(ctx, os.Getenv("TEST_DATABASE_URI"))
	require.NoError(t, err)

	var streamName = testify.StreamName(5)
	var consumerName = testify.ConsumerName(5)
	var topic = testify.Topic(5)

	sub, err := Sub(ctx, pool, clock.New(), streamName, consumerName, topic)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	reportc, _ := sub.Receive(ctx, func(ctx context.Context, events map[string]*entities.StreamEvent) map[string]error {
		var reports = make(map[string]error, len(events))
		for _, event := range events {
			reports[event.EventId] = errors.New(testify.Fake.Emoji().EmojiCode())
		}
		return reports
	})

	// send events
	pub, err := Pub(context.Background(), pool, streamName)
	require.NoError(t, err)

	events := testify.GenStreamEvents(ctx, topic, 1000)
	require.NoError(t, pub.Send(ctx, events))

	reports := <-reportc
	cancel()

	require.Equal(t, DefaultSubscriberSize, len(reports))
}
