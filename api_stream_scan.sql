--->>> api_stream_scan
SELECT id, subject
FROM %s
WHERE id > @cursor
ORDER BY id
LIMIT @size
---<<< api_stream_scan