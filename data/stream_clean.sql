BEGIN;

DO $$
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
END $$;

COMMIT;