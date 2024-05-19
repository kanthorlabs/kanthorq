BEGIN;
CREATE OR REPLACE FUNCTION kanthorq_consumer_pull(cname VARCHAR(128), size SMALLINT)
RETURNS TABLE (cursor_current VARCHAR(128), cursor_next VARCHAR(128)) AS $$
DECLARE 
    consumer RECORD;
    consumer_cursor_next VARCHAR(128);
BEGIN
    -- Select the topic and cursor with a lock
    -- ignore already locked row
    SELECT *
    INTO consumer
    FROM kanthorq_consumer 
    WHERE name = cname
    FOR UPDATE SKIP LOCKED;

    IF consumer.name IS NULL THEN
        RETURN QUERY
        SELECT 'ERROR.CONSUMER.NOT_FOUND', NULL;
    END IF;

    -- Insert new jobs and get the new cursor value
    WITH jobs AS (
        INSERT INTO kanthorq_consumer_job (consumer_name, event_id, tier, topic)
        SELECT cname as consumer_name, event_id, tier, topic
        FROM (
            SELECT DISTINCT ON (event_id) event_id, tier, topic 
            FROM kanthorq_stream
            WHERE tier = consumer.tier AND topic = consumer.topic AND event_id > consumer.cursor 
            ORDER BY event_id
            LIMIT size
        ) sub
        ON CONFLICT(consumer_name, event_id) DO UPDATE 
        SET pull_count = kanthorq_consumer_job.pull_count + 1
        RETURNING event_id
    )
    SELECT max(event_id) INTO consumer_cursor_next FROM jobs;

    IF consumer_cursor_next IS NULL THEN
        RETURN QUERY
        SELECT consumer.cursor, 'ERROR.CONSUMER_JOB.NO_JOB';
    END IF;
    
    INSERT INTO kanthorq_consumer (name, tier, topic, cursor) 
    VALUES(consumer.name, consumer.tier, consumer.topic, consumer_cursor_next)
    ON CONFLICT(name) DO UPDATE 
    SET cursor = EXCLUDED.cursor;

    RETURN QUERY
    SELECT consumer.cursor, consumer_cursor_next;
END;
$$ LANGUAGE plpgsql;

COMMIT;