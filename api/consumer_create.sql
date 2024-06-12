-- >>> consumer_create

select pg_advisory_xact_lock(%d);

CREATE TABLE IF NOT EXISTS %s (
	event_id VARCHAR(64) NOT NULL,
	topic VARCHAR(128) NOT NULL,
	state SMALLINT NOT NULL DEFAULT 0,
	schedule_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000,
	finalized_at BIGINT NOT NULL DEFAULT 0,
	attempt_count SMALLINT NOT NULL DEFAULT 0,
	attempted_at BIGINT NOT NULL DEFAULT 0,
	created_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000,
	updated_at BIGINT NOT NULL DEFAULT 0,
	PRIMARY KEY (event_id)
);

CREATE INDEX IF NOT EXISTS idx_state_scheduler ON %s USING btree("state", "schedule_at");

-- <<< consumer_create
