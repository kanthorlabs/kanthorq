--->>> consumer_register_collection
CREATE TABLE IF NOT EXISTS %s (
	event_id VARCHAR(64) NOT NULL,
	subject VARCHAR(256) NOT NULL,
	state SMALLINT NOT NULL DEFAULT 1,
	schedule_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000,
	finalized_at BIGINT NOT NULL DEFAULT 0,
	attempt_count SMALLINT NOT NULL DEFAULT 0,
	attempted_at BIGINT NOT NULL DEFAULT 0,
	attempted_error jsonb[] NOT NULL DEFAULT '{}',
	metadata JSONB NOT NULL  DEFAULT '{}',
	created_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000,
	updated_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000,
	PRIMARY KEY (event_id)
);

CREATE INDEX IF NOT EXISTS idx_state_scheduling ON %s USING btree("state", "schedule_at", "event_id");
---<<< consumer_register_collection
