SELECT 
  name, 
  stream_name, 
  topic, 
  cursor, 
  attempt_max, 
  created_at, 
  updated_at 
FROM kanthorq_consumer_registry
WHERE name = @consumer_name FOR UPDATE SKIP LOCKED;