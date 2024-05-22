package testify

import "fmt"

func QueryTruncate(table string) string {
	return fmt.Sprintf("TRUNCATE TABLE public.%s CONTINUE IDENTITY RESTRICT;", table)
}

func QueryTruncateConsumer() string {
	return `DO $$
	DECLARE
			rec RECORD;
			drop_table_sql TEXT;
	BEGIN
			-- Loop through each entry in the kanthorq_consumer table
			FOR rec IN SELECT name FROM public.kanthorq_consumer LOOP
					-- Construct the SQL statement to drop the table
					drop_table_sql := 'DROP TABLE IF EXISTS ' || quote_ident(rec.name) || ' CASCADE;';
					-- Execute the drop table statement
					EXECUTE drop_table_sql;
			END LOOP;
			
			-- Delete all entries from the kanthorq_consumer table
			DELETE FROM public.kanthorq_consumer;
	END $$;`
}
