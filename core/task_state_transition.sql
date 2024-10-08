-- >>> task_state_transition
WITH locked_tasks AS (
  SELECT event_id
  FROM %s AS l_tasks
  WHERE
    l_tasks.state = @from_state
    AND l_tasks.attempt_count <= @attempt_max
    AND l_tasks.schedule_at < @attempted_at
  ORDER BY
    l_tasks.state ASC,
    l_tasks.schedule_at ASC
  LIMIT @size
  FOR UPDATE SKIP LOCKED
)
UPDATE %s AS u_tasks
SET
  state = @to_state,
  attempt_count = attempt_count + 1,
  attempted_at = @attempted_at,
  schedule_at = @schedule_at + ((attempt_count ^ 4) * (1 + RANDOM() * 0.2 - 0.1)) * 1000
FROM locked_tasks
WHERE u_tasks.event_id = locked_tasks.event_id 
RETURNING 
  u_tasks.event_id,
  u_tasks.subject,
  u_tasks.state,
  u_tasks.schedule_at,
  u_tasks.attempt_count,
  u_tasks.attempted_at,
  u_tasks.finalized_at,
  u_tasks.created_at,
  u_tasks.updated_at;
---<<< task_state_transition