--->>> api_task_fulfill_from_event
SELECT id, topic, body, metadata, created_at
FROM %s
WHERE id IN (%s);
---<<< api_task_fulfill_from_event