package kanthorq

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/common/idx"
	"github.com/stretchr/testify/require"
)

var uri = "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable"
var table = "kanthorq_stream"
var tier = "default"
var topic = "poc.testing"

func TestConsumerPull(t *testing.T) {
	seed(t)
}

func seed(t *testing.T) {
	conn, err := pgx.Connect(context.Background(), uri)
	require.NoError(t, err)
	defer conn.Close(context.Background())

	var count = 100000000
	var size = 1000000
	var columns = []string{"tier", "topic", "event_id"}
	var entries = make([][]any, size)

	// already seed, ignore
	var found int64
	err = conn.QueryRow(context.Background(), fmt.Sprintf("SELECT COUNT(*) AS found FROM %s", table)).Scan(&found)
	require.NoError(t, err)

	if found == int64(count) {
		return
	}

	// truncate then seed again
	_, err = conn.Exec(context.Background(), fmt.Sprintf("TRUNCATE TABLE public.%s CONTINUE IDENTITY RESTRICT;", table))
	require.NoError(t, err)

	for i := 0; i < count/size; i++ {
		for j := 0; j < size; j++ {
			entries[j] = []any{tier, topic, idx.New("evt")}
		}

		rows, err := conn.CopyFrom(
			context.Background(),
			pgx.Identifier{"kanthorq_stream"},
			columns,
			pgx.CopyFromRows(entries),
		)
		require.NoError(t, err)
		require.Equal(t, int64(size), rows)
	}
}
