--->>> api_consumer_lock
SELECT 
  stream_id, 
  stream_name, 
  id, 
  name, 
  subject_filter, 
  cursor, 
  attempt_max, 
  created_at, 
  updated_at 
FROM kanthorq_consumer_registry
WHERE name = @consumer_name FOR UPDATE SKIP LOCKED;
---<<< api_consumer_lock