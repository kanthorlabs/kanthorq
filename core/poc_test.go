package core

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/common/idx"
	"github.com/stretchr/testify/require"
)

var uri = "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable"
var tier = "default"
var topic = "poc.testing"
var consumer = "poc"

func TestPOC(t *testing.T) {
	conn, err := pgx.Connect(context.Background(), uri)
	require.NoError(t, err)
	defer conn.Close(context.Background())

	seed(t, conn)

	var cursor string
	err = conn.
		QueryRow(
			context.Background(),
			QueryConsumerPull,
			pgx.NamedArgs{
				"consumer_name": consumer,
				"size":          100,
			},
		).
		Scan(&cursor)
	require.NoError(t, err)

	require.NotEmpty(t, cursor)
}

func seed(t *testing.T, conn *pgx.Conn) {
	var err error

	var count = int64(100000000)
	var size = int64(1000000)
	var entries = make([][]any, size)

	// already seed, ignore
	var found int64
	err = conn.QueryRow(context.Background(), fmt.Sprintf("SELECT COUNT(*) AS found FROM %s", CollectionStream)).Scan(&found)
	require.NoError(t, err)

	if found >= count {
		return
	}

	for i := int64(0); i < (count-found)/size; i++ {
		for j := int64(0); j < size; j++ {
			entries[j] = []any{tier, topic, idx.New("evt")}
		}

		rows, err := conn.CopyFrom(
			context.Background(),
			pgx.Identifier{CollectionStream},
			(&Stream{}).Properties(),
			pgx.CopyFromRows(entries),
		)
		require.NoError(t, err)
		require.Equal(t, int64(size), rows)
	}

	_, err = conn.Exec(context.Background(), fmt.Sprintf("TRUNCATE TABLE public.%s CONTINUE IDENTITY RESTRICT;", CollectionConsumerJob))
	require.NoError(t, err)

	_, err = conn.Exec(context.Background(), fmt.Sprintf("TRUNCATE TABLE public.%s CONTINUE IDENTITY RESTRICT;", CollectionConsumer))
	require.NoError(t, err)

	_, err = conn.Exec(
		context.Background(),
		fmt.Sprintf("INSERT INTO public.%s (name, tier, topic, cursor) VALUES (@name, @tier, @topic, '');", CollectionConsumer),
		pgx.NamedArgs{
			"name":  consumer,
			"tier":  tier,
			"topic": topic,
		},
	)
	require.NoError(t, err)
}
