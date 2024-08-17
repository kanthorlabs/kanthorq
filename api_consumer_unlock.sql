--->>> api_consumer_unlock
UPDATE kanthorq_consumer_registry
SET cursor = @consumer_cursor 
WHERE name = @consumer_name
RETURNING stream_id, stream_name, id, name, subject_filter, cursor, attempt_max, created_at, updated_at;
---<<< api_consumer_unlock