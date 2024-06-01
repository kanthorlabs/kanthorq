package kanthorq

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConnection(t *testing.T) {
	pool, err := Connection(context.Background(), os.Getenv("TEST_DATABASE_URI"))
	require.NoError(t, err)
	require.NotNil(t, pool)
}
