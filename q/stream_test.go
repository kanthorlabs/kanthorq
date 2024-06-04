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

func TestStream(t *testing.T) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("TEST_DATABASE_URI"))
	require.NoError(t, err)

	stream, err := Stream(context.Background(), conn, &entities.Stream{Name: testify.StreamName(5)})
	require.NoError(t, err)
	require.NotNil(t, stream)
}
