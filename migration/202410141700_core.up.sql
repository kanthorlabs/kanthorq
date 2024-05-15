BEGIN;

CREATE TABLE IF NOT EXISTS kanthorq_stream_message (
  stream_topic VARCHAR(128) NOT NULL,
  stream_name VARCHAR(128) NOT NULL,
  message_id VARCHAR(128) NOT NULL,
  PRIMARY KEY (stream_name, message_id)
);

CREATE TABLE IF NOT EXISTS kanthorq_consumer (
  name VARCHAR(128) NOT NULL,
  stream_topic VARCHAR(128) NOT NULL,
  stream_cursor VARCHAR(128) NOT NULL,
  PRIMARY KEY (name)
);

CREATE TABLE IF NOT EXISTS kanthorq_consumer_message (
  stream_topic VARCHAR(128) NOT NULL,
  stream_name VARCHAR(128) NOT NULL,
  message_id VARCHAR(128) NOT NULL,
  PRIMARY KEY (stream_name, message_id)
);

COMMIT;