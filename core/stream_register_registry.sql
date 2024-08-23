--->>> stream_register_registry
INSERT INTO kanthorq_stream_registry(id, name)
VALUES(@stream_id, @stream_name)
ON CONFLICT(name) DO UPDATE 
SET updated_at = EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000
RETURNING id, name, created_at, updated_at;
---<<< stream_register_registry