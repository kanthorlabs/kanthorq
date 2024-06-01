-- stream_ensure
INSERT INTO kanthorq_stream(name)
VALUES(@stream_name)
ON CONFLICT(name) DO UPDATE 
SET updated_at = EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000
RETURNING name, created_at, updated_at;