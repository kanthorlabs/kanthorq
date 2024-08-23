--->>> stream_get_event
SELECT id, subject, body, metadata, created_at
FROM %s
WHERE id IN (%s);
---<<< stream_get_event