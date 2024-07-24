--->>> api_task_convert_from_event
INSERT INTO %s (event_id, topic, state)
SELECT event_id, topic, @intial_state
FROM %s
WHERE topic = @consumer_topic AND event_id > @consumer_cursor
ORDER BY topic, event_id
LIMIT @size
ON CONFLICT(event_id) DO UPDATE
SET updated_at = EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000
RETURNING event_id
---<<< api_task_convert_from_event