BEGIN;
  
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'stream_event') THEN
  CREATE TYPE stream_event AS (
      topic varchar(128),
      event_id varchar(64),
      created_at BIGINT 
  );
  END IF;
END $$;

CREATE TABLE IF NOT EXISTS kanthorq_stream (
  name VARCHAR(128) NOT NULL,
  created_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000,
  updated_at BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (name)
);

CREATE TABLE IF NOT EXISTS kanthorq_consumer (
  name VARCHAR(128) NOT NULL,
  stream_name VARCHAR(128) NOT NULL,
  topic VARCHAR(128) NOT NULL,
  cursor VARCHAR(64) NOT NULL,
  created_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000,
  updated_at BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (name)
);

COMMIT;
