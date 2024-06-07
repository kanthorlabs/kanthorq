package q

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/stretchr/testify/require"
)

func TestConsumer(t *testing.T) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("KANTHORQ_POSTGRES_URI"))
	require.NoError(t, err)

	c := &entities.Consumer{
		StreamName: testify.StreamName(5),
		Name:       testify.ConsumerName(5),
		Topic:      testify.Topic(5),
	}
	consuemr, err := Consumer(ctx, conn, c)
	require.NoError(t, err)
	require.NotNil(t, consuemr)
}
