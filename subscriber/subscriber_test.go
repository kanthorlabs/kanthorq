package subscriber

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/kanthorlabs/kanthorq"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/publisher"
	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/stretchr/testify/require"
)

func TestSubscriber(t *testing.T) {
	ctx := context.Background()

	pool, err := kanthorq.Connection(ctx, os.Getenv("TEST_DATABASE_URI"))
	require.NoError(t, err)

	var streamName = testify.StreamName(5)
	var consumerName = testify.ConsumerName(5)
	var topic = testify.Topic(5)

	// publish first
	pub := publisher.New(&publisher.Config{StreamName: streamName}, pool)
	require.NoError(t, pub.Start(ctx))
	defer require.NoError(t, pub.Stop(ctx))

	events := testify.GenStreamEvents(ctx, topic, 1000)
	require.NoError(t, pub.Send(ctx, events))

	// then subscribe
	conf := &Config{
		StreamName:   streamName,
		ConsumerName: consumerName,
		Topic:        topic,
	}
	sub := New(conf, pool)
	require.NoError(t, sub.Start(ctx))
	// don't use defer here to stop the subscriber
	// otherwise the select state could not receive data

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

	require.NoError(t, sub.Stop(ctx))
}
