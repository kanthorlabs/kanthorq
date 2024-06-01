-- consumer_job_pull
SELECT topic, event_id, created_at
FROM %s AS stream
WHERE (topic, event_id) IN (%s)
ORDER BY topic, event_id;