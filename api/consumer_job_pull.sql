-- consumer_job_pull
SELECT topic, event_id, body, metadata, created_at
FROM %s AS stream
WHERE event_id IN (%s)
ORDER BY event_id;