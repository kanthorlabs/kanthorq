--->>> api_stream_register_collection
CREATE TABLE IF NOT EXISTS %s (
	id VARCHAR(64) NOT NULL,
	subject VARCHAR(128) NOT NULL,
  body BYTEA NOT NULL,
	metadata jsonb NOT NULL,
	created_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000,
	PRIMARY KEY (id)
);
---<<< api_stream_register_collection
