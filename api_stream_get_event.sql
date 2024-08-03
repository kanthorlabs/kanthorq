--->>> api_stream_get_event
SELECT id, subject, body, metadata, created_at
FROM %s
WHERE id IN (%s);
---<<< api_stream_get_event