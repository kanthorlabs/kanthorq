package kanthorq

import (
	"context"
	"os"
	"testing"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/stretchr/testify/require"
)

func TestStream(t *testing.T) {
	pool, err := Connection(context.Background(), os.Getenv("TEST_DATABASE_URI"))
	require.NoError(t, err)
	require.NotNil(t, pool)

	stream, err := Stream(context.Background(), pool, &entities.Stream{Name: testify.StreamName(5)})
	require.NoError(t, err)
	require.NotNil(t, stream)
}
