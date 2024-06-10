-- >>> stream_create

CREATE TABLE IF NOT EXISTS %s (
	event_id VARCHAR(64) NOT NULL,
	topic VARCHAR(128) NOT NULL,
  body jsonb NOT NULL,
	metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
	created_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000,
	PRIMARY KEY (event_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_event_id ON %s USING btree("topic", "event_id");

-- <<< stream_create
