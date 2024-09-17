--->>> consumer_register_registry
INSERT INTO kanthorq_consumer_registry(
  stream_id,
  stream_name,
  id,
  name,
  subject_includes,
  subject_excludes,
  cursor,
  attempt_max,
  visibility_timeout
)
VALUES(
  @stream_id,
  @stream_name,
  @consumer_id,
  @consumer_name,
  @consumer_subject_includes,
  @consumer_subject_excludes,
  @consumer_cursor,
  @consumer_attempt_max,
  @consumer_visibility_timeout
)
ON CONFLICT(name) DO UPDATE 
SET updated_at = EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000
RETURNING 
  stream_id,
  stream_name,
  id,
  name,
  kind,
  subject_includes,
  subject_excludes,
  cursor,
  attempt_max,
  visibility_timeout,
  created_at,
  updated_at;
---<<< consumer_register_registry