package kanthorq

import (
	"context"
	"os"
	"testing"

	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/stretchr/testify/require"
)

func TestPublisher(t *testing.T) {
	pool, err := Connection(context.Background(), os.Getenv("TEST_DATABASE_URI"))
	require.NoError(t, err)

	name := testify.StreamName(5)

	pub, err := Pub(context.Background(), pool, name)
	require.NoError(t, err)

	events := testify.GenStreamEvents(context.Background(), testify.Topic(5), 1000)

	err = pub.Send(context.Background(), events)
	require.NoError(t, err)
}
