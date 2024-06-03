package publisher

import (
	"context"
	"os"
	"testing"

	"github.com/kanthorlabs/kanthorq"
	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/stretchr/testify/require"
)

func TestPublisher(t *testing.T) {
	ctx := context.Background()

	pool, err := kanthorq.Connection(ctx, os.Getenv("TEST_DATABASE_URI"))
	require.NoError(t, err)

	conf := &Config{StreamName: testify.StreamName(5)}
	pub := New(conf, pool)
	require.NoError(t, pub.Start(ctx))
	defer require.NoError(t, pub.Stop(ctx))

	events := testify.GenStreamEvents(ctx, testify.Topic(5), 1000)
	err = pub.Send(ctx, events)
	require.NoError(t, err)
}
