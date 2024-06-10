-- >>> consumer_pull
WITH jobs AS (
  INSERT INTO %s (event_id, topic)
  SELECT event_id, topic
  FROM %s
  WHERE topic = @consumer_topic AND event_id > @consumer_cursor
  ORDER BY topic, event_id
  LIMIT @size
  ON CONFLICT(event_id) DO UPDATE
  SET updated_at = EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000
  RETURNING event_id
),
next_cursor AS (
	SELECT max(event_id) AS next_event FROM jobs
)
UPDATE kanthorq_consumer
SET cursor = next_event
FROM next_cursor
WHERE name = @consumer_name AND next_event IS NOT NULL
RETURNING cursor;	
-- <<< consumer_pull
