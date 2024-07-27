--->>> api_consumer_register_registry
INSERT INTO kanthorq_consumer_registry(stream_id, stream_name, id, name, topic, attempt_max)
VALUES(@stream_id, @stream_name, @consumer_id, @consumer_name, @consumer_topic, @consumer_attempt_max)
ON CONFLICT(name) DO UPDATE 
SET updated_at = EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000
RETURNING stream_id, stream_name, id, name, topic, cursor, attempt_max, created_at, updated_at;
---<<< api_consumer_register_registry