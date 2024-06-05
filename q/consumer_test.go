package q

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/stretchr/testify/require"
)

func TestConsumer(t *testing.T) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, os.Getenv("TEST_POSTGRES_URI"))
	require.NoError(t, err)

	c := &entities.Consumer{
		StreamName: testify.StreamName(5),
		Name:       testify.ConsumerName(5),
		Topic:      testify.Topic(5),
	}
	consuemr, err := Consumer(ctx, pool, c)
	require.NoError(t, err)
	require.NotNil(t, consuemr)
}
