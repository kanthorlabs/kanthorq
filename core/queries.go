package core

import "github.com/jackc/pgx/v5"

func QueryConsumerEnsure(name, topic string) (string, pgx.NamedArgs) {
	statement := `SELECT name, topic, cursor, created_at, updated_at FROM consumer_ensure(@consumer_name, @topic);`
	args := pgx.NamedArgs{
		"consumer_name": name,
		"topic":         topic,
	}
	return statement, args
}

func QueryConsumerPull(name string, size int) (string, pgx.NamedArgs) {
	statement := `SELECT consumer_name, cursor_current, cursor_next FROM kanthorq_consumer_pull(@consumer_name, CAST(@size as SMALLINT));`
	args := pgx.NamedArgs{
		"consumer_name": name,
		"size":          size,
	}
	return statement, args
}
