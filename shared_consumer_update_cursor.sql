UPDATE kanthorq_consumer_registry
SET cursor = @consumer_cursor 
WHERE name = @consumer_name
RETURNING stream_id, stream_name, id, name, topic, cursor, attempt_max, created_at, updated_at;