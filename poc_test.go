package kanthorq

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/common/idx"
	"github.com/stretchr/testify/require"
)

var stream = "poc"
var topic = "poc.testing"

func TestConsumerPull(t *testing.T) {
	var uri = "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable"
	conn, err := pgx.Connect(context.Background(), uri)
	require.NoError(t, err)
	defer conn.Close(context.Background())

	seed(t, conn)
}

func seed(t *testing.T, conn *pgx.Conn) {
	var columns = []string{"stream_topic", "stream_name", "message_id"}

	var count = gofakeit.Number(900000, 1000000)
	var entries = make([][]any, count)
	for i := 0; i < count; i++ {
		entries[i] = []any{topic, stream, idx.New("msg")}
	}

	_, err := conn.Exec(context.Background(), "TRUNCATE TABLE public.kanthorq_stream_message CONTINUE IDENTITY RESTRICT;")
	require.NoError(t, err)

	rows, err := conn.CopyFrom(
		context.Background(),
		pgx.Identifier{"kanthorq_stream_message"},
		columns,
		pgx.CopyFromRows(entries),
	)
	require.NoError(t, err)
	require.Equal(t, int64(count), rows)
}
