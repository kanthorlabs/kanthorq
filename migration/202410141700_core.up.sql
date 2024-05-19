BEGIN;

CREATE TABLE IF NOT EXISTS kanthorq_stream (
  tier VARCHAR(128) NOT NULL,
  topic VARCHAR(128) NOT NULL,
  event_id VARCHAR(64) NOT NULL,
  PRIMARY KEY (tier, topic, event_id)
);

CREATE TABLE IF NOT EXISTS kanthorq_consumer (
  name VARCHAR(128) NOT NULL,
  tier VARCHAR(128) NOT NULL,
  topic VARCHAR(128) NOT NULL,
  cursor VARCHAR(128) NOT NULL,
  PRIMARY KEY (name)
);

CREATE TABLE IF NOT EXISTS kanthorq_consumer_job (
  consumer_name VARCHAR(128) NOT NULL,
  event_id VARCHAR(64) NOT NULL,
  tier VARCHAR(128) NOT NULL,
  topic VARCHAR(128) NOT NULL,
  pull_count SMALLINT NOT NULL DEFAULT 1, 
  PRIMARY KEY (consumer_name, event_id)
);

COMMIT;