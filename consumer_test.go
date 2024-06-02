package kanthorq

import (
	"context"
	"os"
	"testing"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/stretchr/testify/require"
)

func TestConsumer(t *testing.T) {
	pool, err := Connection(context.Background(), os.Getenv("TEST_DATABASE_URI"))
	require.NoError(t, err)
	require.NotNil(t, pool)

	c := &entities.Consumer{
		StreamName: testify.StreamName(5),
		Name:       testify.ConsumerName(5),
		Topic:      testify.Topic(5),
	}
	consuemr, err := Consumer(context.Background(), pool, c)
	require.NoError(t, err)
	require.NotNil(t, consuemr)
}
