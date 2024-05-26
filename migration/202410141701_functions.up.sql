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
      CREATE UNIQUE INDEX IF NOT EXISTS idx_event_id ON "kanthorq_stream_%s" USING btree("event_id");
    $QUERY$,
    req_stream_name,
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
BEGIN
  -- the default configuration of PostgreSQL is to treat all unquoted identifiers (such as table names) as case-insensitive
  -- make sure all table names are quoted

  INSERT INTO kanthorq_consumer(name, stream_name, topic, cursor)
  VALUES(req_consumer_name, req_stream_name, req_topic, '')
  ON CONFLICT(name) DO UPDATE 
  SET updated_at = EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000
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
        topic VARCHAR(128) NOT NULL,
        state SMALLINT NOT NULL DEFAULT 0,
        schedule_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000,
        attempt_count SMALLINT NOT NULL DEFAULT 0,
        attempted_at BIGINT NOT NULL DEFAULT 0,
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

COMMIT;