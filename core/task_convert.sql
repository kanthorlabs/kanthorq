--->>> task_convert
INSERT INTO %s (event_id, subject, state)
SELECT id, subject, @intial_state::SMALLINT as state
FROM %s
WHERE id IN (%s)
ON CONFLICT(event_id) DO NOTHING
RETURNING 
  event_id,
  subject,
  state,
  schedule_at,
  attempt_count,
  attempted_at,
  finalized_at,
  created_at,
  updated_at;
---<<< task_convert