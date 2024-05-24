package testify

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func QueryTruncateConsumer() func(ctx context.Context, conn *pgxpool.Pool) error {
	statement := `DO $$
	DECLARE
			rec RECORD;
			drop_table_sql TEXT;
	BEGIN
			-- Loop through each entry in the kanthorq_consumer table
			FOR rec IN SELECT name FROM kanthorq_consumer LOOP
					-- Construct the SQL statement to drop the table
					drop_table_sql := 'DROP TABLE IF EXISTS "kanthorq_consumer_' || rec.name || '" CASCADE;';
					-- Execute the drop table statement
					EXECUTE drop_table_sql;
			END LOOP;
			
			-- Delete all entries from the kanthorq_consumer table
			DELETE FROM kanthorq_consumer;
	END $$;`

	return func(ctx context.Context, conn *pgxpool.Pool) error {
		_, err := conn.Exec(ctx, statement)
		return err
	}
}

func QueryTruncateStream() func(ctx context.Context, conn *pgxpool.Pool) error {
	statement := `DO $$
	DECLARE
			rec RECORD;
			drop_table_sql TEXT;
	BEGIN
			-- Loop through each entry in the kanthorq_stream table
			FOR rec IN SELECT name FROM kanthorq_stream LOOP
					-- Construct the SQL statement to drop the table
					drop_table_sql := 'DROP TABLE IF EXISTS "kanthorq_stream_' || rec.name || '" CASCADE;';
					-- Execute the drop table statement
					EXECUTE drop_table_sql;
			END LOOP;
			
			-- Delete all entries from the kanthorq_stream table
			DELETE FROM kanthorq_stream;
	END $$;`

	return func(ctx context.Context, conn *pgxpool.Pool) error {
		_, err := conn.Exec(ctx, statement)
		return err
	}
}
