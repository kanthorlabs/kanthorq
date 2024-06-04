package publisher

import (
	"context"
	"os"
	"testing"

	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/stretchr/testify/require"
)

func TestPublisher(t *testing.T) {
	ctx := context.Background()

	pub, err := New(ctx, &Config{
		ConnectionUri: os.Getenv("TEST_DATABASE_URI"),
		StreamName:    testify.StreamName(5),
	})
	require.NoError(t, err)

	events := testify.GenStreamEvents(ctx, testify.Topic(5), 1000)
	require.NoError(t, pub.Send(ctx, events))
}
