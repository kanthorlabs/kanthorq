--->>> stream_scan
SELECT id, subject
FROM %s
WHERE id > @cursor
ORDER BY id
LIMIT @size
---<<< stream_scan