--->>> consumer_ensure
INSERT INTO kanthorq_consumer(name, stream_name, topic, cursor)
VALUES(@consumer_name, @stream_name, @topic, '')
ON CONFLICT(name) DO UPDATE 
SET updated_at = EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000
RETURNING name, stream_name, topic, cursor, created_at, updated_at;
---<<< consumer_ensure