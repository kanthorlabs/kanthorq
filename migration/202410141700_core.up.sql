BEGIN;

CREATE TABLE IF NOT EXISTS kanthorq_task (
  org_id VARCHAR(128) NOT NULL,
  id VARCHAR(128) NOT NULL,
  created_at BIGINT NOT NULL,
  updated_at BIGINT NOT NULL DEFAULT 0,

  -- The message read counter keeps track of how many times the message has been read.
  read_couter SMALLINT NOT NULL,
  -- The UTC timestamp indicates that the message can only be read after that time.
  read_after BIGINT NOT NULL,
  -- The UTC timestamp indicates when the message was read by a worker.
  read_at BIGINT NOT NULL,
  -- The column indicates which worker read the message.
  read_by VARCHAR(128) NOT NULL,

  message JSONB,

  PRIMARY KEY (org_id, id)
);

COMMIT;