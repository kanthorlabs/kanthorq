-- >>> api_stream_register_collection

CREATE TABLE IF NOT EXISTS %s (
	id VARCHAR(64) NOT NULL,
	topic VARCHAR(128) NOT NULL,
  body BYTEA NOT NULL,
	metadata jsonb NOT NULL,
	created_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000,
	PRIMARY KEY (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_topic_sharding ON %s USING btree("topic", "id");

-- <<< api_stream_register_collection
