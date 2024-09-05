--->>> consumer_lock
SELECT 
  stream_id,
  stream_name,
  id,
  name,
  subject_includes,
  subject_excludes,
  cursor,
  attempt_max,
  visibility_timeout,
  created_at,
  updated_at 
FROM kanthorq_consumer_registry
WHERE name = @consumer_name FOR UPDATE SKIP LOCKED;
---<<< consumer_lock