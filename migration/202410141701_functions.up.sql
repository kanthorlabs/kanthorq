BEGIN;

CREATE OR REPLACE FUNCTION kanthorq_consumer_pull(consumer_name VARCHAR(128), size SMALLINT) RETURNS VARCHAR(128) AS $$
DECLARE 
    consumer RECORD;
    consumer_cursor_next VARCHAR(128);
BEGIN
    -- Select the topic and cursor with a lock
    -- ignore already locked row
    SELECT *
    INTO consumer
    FROM kanthorq_consumer 
    WHERE name = consumer_name
    FOR UPDATE SKIP LOCKED;

    -- If the select statement from kanthorq_consumer returned nothing
    -- then consumer.topic and consumer.cursor will be null
    -- A valid consumer must have not null and not empty consumer
    -- but they may have empty cursor when it was initialized
    IF consumer.tier IS NOT NULL THEN
        
        -- Insert new jobs and get the new cursor value
        WITH jobs AS (
            INSERT INTO kanthorq_consumer_job (tier, topic, event_id)
            SELECT tier, topic, event_id 
            FROM kanthorq_stream
            WHERE tier = consumer.tier AND topic = consumer.topic AND event_id > consumer.cursor 
            ORDER BY tier ASC, topic ASC, event_id ASC
            LIMIT size
            ON CONFLICT(tier, topic, event_id) DO UPDATE 
            SET write_count = kanthorq_consumer_job.write_count + 1
            RETURNING event_id
        )
        SELECT max(event_id) INTO consumer_cursor_next FROM jobs;

        IF consumer_cursor_next IS NOT NULL AND consumer_cursor_next <> '' THEN
            -- Update the cursor in kanthorq_consumer
            INSERT INTO kanthorq_consumer (name, tier, topic, cursor) 
            VALUES(consumer.name, consumer.tier, consumer.topic, consumer_cursor_next)
            ON CONFLICT(name) DO UPDATE 
            SET cursor = EXCLUDED.cursor;
        END IF;
    END IF;

    -- Return the new cursor value
    RETURN consumer_cursor_next;
END;
$$ LANGUAGE plpgsql;

COMMIT;