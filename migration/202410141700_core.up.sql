BEGIN;

CREATE TABLE IF NOT EXISTS kanthorq_stream_registry (
  id VARCHAR(64) NOT NULL,
  name VARCHAR(256) NOT NULL,
  created_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000,
  updated_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000,
  PRIMARY KEY (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_stream_name ON kanthorq_stream_registry USING btree("name");

CREATE TABLE IF NOT EXISTS kanthorq_consumer_registry (
  stream_id VARCHAR(64) NOT NULL,
  stream_name VARCHAR(256) NOT NULL,
  id VARCHAR(64) NOT NULL,
  name VARCHAR(256) NOT NULL,
  kind SMALLINT NOT NULL DEFAULT 1,
  subject_includes VARCHAR(256) ARRAY NOT NULL,
  subject_excludes VARCHAR(256) ARRAY NOT NULL DEFAULT '{}',
  cursor VARCHAR(64) NOT NULL,
  attempt_max SMALLINT NOT NULL DEFAULT 16,
  visibility_timeout BIGINT NOT NULL DEFAULT 300000,
  created_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000,
  updated_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000,
  PRIMARY KEY (name)
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_consumer_name ON kanthorq_consumer_registry USING btree("name");

COMMIT;
