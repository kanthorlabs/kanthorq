package q

import (
	"context"
	"testing"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/stretchr/testify/require"
)

func TestConsumer(t *testing.T) {
	ctx := context.Background()
	conn, err := testify.SetupPostgres(ctx)
	require.NoError(t, err)
	defer conn.Close(ctx)

	c := &entities.Consumer{
		StreamName: testify.StreamName(5),
		Name:       testify.ConsumerName(5),
		Topic:      testify.Topic(5),
	}
	consumer, err := NewConsumer(ctx, conn, c)
	require.NoError(t, err)
	require.NotNil(t, consumer)
}
