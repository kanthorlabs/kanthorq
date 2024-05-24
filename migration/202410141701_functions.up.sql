BEGIN;

CREATE OR REPLACE FUNCTION stream_ensure(req_stream_name VARCHAR(128))
RETURNS kanthorq_stream AS $$
DECLARE 
    stream kanthorq_stream;
    stream_create_sql TEXT;
    ts BIGINT := EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000;
BEGIN
    -- the default configuration of PostgreSQL is to treat all unquoted identifiers (such as table names) as case-insensitive
    -- make sure all table names are quoted

    INSERT INTO kanthorq_stream(name)
    VALUES(req_stream_name)
    ON CONFLICT(name) DO UPDATE 
    SET updated_at = ts
    RETURNING * INTO stream;

    stream_create_sql := FORMAT(
        $QUERY$
        CREATE TABLE IF NOT EXISTS "kanthorq_stream_%s" (
            topic VARCHAR(128) NOT NULL,
            event_id VARCHAR(64) NOT NULL,
            created_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000,
            PRIMARY KEY (topic, event_id)
        );
        $QUERY$,
        req_stream_name
    );
    EXECUTE stream_create_sql;

    RETURN stream;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION consumer_ensure(req_consumer_name VARCHAR(128), req_stream_name VARCHAR(128), req_topic VARCHAR(128))
RETURNS kanthorq_consumer AS $$
DECLARE 
    consumer kanthorq_consumer;
    consumer_create_sql TEXT;
    ts BIGINT := EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000;
BEGIN
    -- the default configuration of PostgreSQL is to treat all unquoted identifiers (such as table names) as case-insensitive
    -- make sure all table names are quoted

    INSERT INTO kanthorq_consumer(name, stream_name, topic, cursor)
    VALUES(req_consumer_name, req_stream_name, req_topic, '')
    ON CONFLICT(name) DO UPDATE 
    SET updated_at = ts
    RETURNING * INTO consumer;

    -- if the request topic is not matched with the existing topic
    -- there is something wrong in the request
    IF consumer.topic <> req_topic THEN
        RAISE EXCEPTION 'ERROR.CONSUMER.REQUEST_TOPIC.NOT_MATCH: EXPECTED:% ACTUAL:%', consumer.name, req_topic;
    END IF;

    consumer_create_sql := FORMAT(
        $QUERY$
        CREATE TABLE IF NOT EXISTS "kanthorq_consumer_%s" (
            event_id VARCHAR(64) NOT NULL,
            name VARCHAR(128) NOT NULL,
            topic VARCHAR(128) NOT NULL,
            state SMALLINT NOT NULL DEFAULT 0,
            pull_count SMALLINT NOT NULL DEFAULT 0,
            created_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000,
            updated_at BIGINT NOT NULL DEFAULT 0,
            PRIMARY KEY (event_id)
        );
        $QUERY$,
        req_consumer_name
    );
    EXECUTE consumer_create_sql;

    RETURN consumer;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION kanthorq_consumer_pull(req_consumer_name VARCHAR(128), req_pull_size SMALLINT)
RETURNS TABLE (consumer_name VARCHAR(128), cursor_current VARCHAR(64), cursor_next VARCHAR(64)) AS $$
DECLARE 
    consumer RECORD;
    consumer_cursor_next VARCHAR(64);
    consumer_job_insert_sql TEXT;
BEGIN
    -- the default configuration of PostgreSQL is to treat all unquoted identifiers (such as table names) as case-insensitive
    -- make sure all table names are quoted

    -- Select the topic and cursor with a lock
    -- ignore already locked row
    SELECT *
    INTO consumer
    FROM kanthorq_consumer AS kqc
    WHERE kqc.name = req_consumer_name
    FOR UPDATE SKIP LOCKED;

    IF consumer.name IS NULL THEN
        RAISE EXCEPTION 'ERROR.CONSUMER.BUSY: %', consumer.name;
    END IF;

    -- Insert new jobs and get the new cursor value
    consumer_job_insert_sql := FORMAT(
        $QUERY$
        WITH jobs AS (
            INSERT INTO "kanthorq_consumer_%s" (event_id, name, topic)
                SELECT event_id, %L as name, topic
                FROM "kanthorq_stream_%s"
                WHERE topic = %L AND event_id > %L 
                ORDER BY event_id
                LIMIT %s
            ON CONFLICT(event_id) DO UPDATE 
            SET pull_count = "kanthorq_consumer_%s".pull_count + 1
            RETURNING event_id
        )
        SELECT max(event_id) AS consumer_cursor_next FROM jobs;
        $QUERY$,
        consumer.name,
        consumer.name,
        consumer.stream_name,
        consumer.topic,
        consumer.cursor,
        req_pull_size,
        consumer.name
    );
    EXECUTE consumer_job_insert_sql INTO consumer_cursor_next;

    IF consumer_cursor_next IS NOT NULL THEN
        UPDATE kanthorq_consumer SET cursor = consumer_cursor_next WHERE name = consumer.name;
    END IF;

    -- we should return all NULL or all STRING
    -- consumer_cursor_next maybe null because of no more job
    -- should cast it as STRING if it is NULL
    RETURN QUERY SELECT consumer.name, consumer.cursor, COALESCE(consumer_cursor_next, '');
END;
$$ LANGUAGE plpgsql;

COMMIT;