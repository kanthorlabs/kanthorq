--->>> api_task_convert_from_event
INSERT INTO %s (event_id, topic, state)
SELECT id, topic, @intial_state::SMALLINT as state
FROM %s
WHERE 1=1 
  %s
  AND topic LIKE @consumer_topic_filter
ORDER BY id
LIMIT @size
ON CONFLICT(event_id) DO NOTHING
RETURNING event_id, topic, state, schedule_at, attempt_count, attempted_at, finalized_at, created_at, updated_at;
---<<< api_task_convert_from_event