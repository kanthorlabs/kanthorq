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

func TestStream(t *testing.T) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, os.Getenv("KANTHORQ_POSTGRES_URI"))
	require.NoError(t, err)

	stream, err := Stream(context.Background(), pool, &entities.Stream{Name: testify.StreamName(5)})
	require.NoError(t, err)
	require.NotNil(t, stream)
}
