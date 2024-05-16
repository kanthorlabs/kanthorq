BEGIN;

CREATE TABLE IF NOT EXISTS kanthorq_stream (
  tier VARCHAR(128) NOT NULL,
  topic VARCHAR(128) NOT NULL,
  event_id VARCHAR(128) NOT NULL,
  PRIMARY KEY (tier, topic, event_id)
);

CREATE TABLE IF NOT EXISTS kanthorq_consumer (
  name VARCHAR(128) NOT NULL,
  topic VARCHAR(128) NOT NULL,
  cursor VARCHAR(128) NOT NULL,
  PRIMARY KEY (name)
);

CREATE TABLE IF NOT EXISTS kanthorq_consumer_job (
  tier VARCHAR(128) NOT NULL,
  topic VARCHAR(128) NOT NULL,
  event_id VARCHAR(128) NOT NULL,
  PRIMARY KEY (tier, topic, event_id)
);

COMMIT;