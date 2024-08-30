--->>> task_convert
INSERT INTO %s (event_id, subject, state, schedule_at, metadata)
SELECT id, subject, @intial_state::SMALLINT as state, @schedule_at as schedule_at, metadata
FROM %s
WHERE id IN (%s)
ON CONFLICT(event_id) DO NOTHING
RETURNING 
  event_id,
  subject,
  state,
  schedule_at,
  finalized_at,
  attempt_count,
  attempted_at,
  attempted_error,
  metadata,
  created_at,
  updated_at;
---<<< task_convert