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

	pub := New(&Config{
		ConnectionUri: os.Getenv("TEST_POSTGRES_URI"),
		StreamName:    testify.StreamName(5),
	})
	require.NoError(t, pub.Start(ctx))
	defer func() {
		require.NoError(t, pub.Stop(ctx))
	}()

	events := testify.GenStreamEvents(ctx, testify.Topic(5), 1000)
	require.NoError(t, pub.Send(ctx, events))
}
