package q

import (
	"context"
	"testing"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/stretchr/testify/require"
)

func TestStream(t *testing.T) {
	ctx := context.Background()
	conn, err := testify.SetupPostgres(ctx)
	require.NoError(t, err)
	defer conn.Close(ctx)

	stream, err := NewStream(context.Background(), conn, &entities.Stream{
		Name: testify.StreamName(5),
	})
	require.NoError(t, err)
	require.NotNil(t, stream)
}
