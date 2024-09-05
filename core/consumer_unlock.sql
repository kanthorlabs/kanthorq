--->>> consumer_unlock
UPDATE kanthorq_consumer_registry
SET cursor = @consumer_cursor 
WHERE name = @consumer_name
RETURNING 
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
  updated_at;
---<<< consumer_unlock