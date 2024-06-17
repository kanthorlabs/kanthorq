BEGIN;

CREATE TABLE IF NOT EXISTS kanthorq_stream (
  name VARCHAR(128) NOT NULL,
  created_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000,
  updated_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000,
  PRIMARY KEY (name)
);

CREATE TABLE IF NOT EXISTS kanthorq_consumer (
  name VARCHAR(128) NOT NULL,
  stream_name VARCHAR(128) NOT NULL,
  topic VARCHAR(128) NOT NULL,
  cursor VARCHAR(64) NOT NULL,
  created_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000,
  updated_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000,
  PRIMARY KEY (name)
);

COMMIT;
