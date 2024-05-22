BEGIN;

CREATE TABLE IF NOT EXISTS kanthorq_stream (
  topic VARCHAR(128) NOT NULL,
  event_id VARCHAR(64) NOT NULL,
  PRIMARY KEY (topic, event_id)
);

CREATE TABLE IF NOT EXISTS kanthorq_consumer (
  name VARCHAR(128) NOT NULL,
  topic VARCHAR(128) NOT NULL,
  cursor VARCHAR(64) NOT NULL,
  PRIMARY KEY (name)
);

COMMIT;