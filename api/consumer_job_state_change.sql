-- consumer_job_state_change
WITH locked_jobs AS (
  SELECT
    event_id
  FROM %s AS l_jobs
  WHERE
    l_jobs.state = @from_state
    AND l_jobs.schedule_at < @attempt_at
  ORDER BY
    l_jobs.state ASC,
    l_jobs.schedule_at ASC
  LIMIT @size
  FOR UPDATE SKIP LOCKED
)
UPDATE %s AS u_jobs
SET
  state = @to_state,
  attempt_count = attempt_count + 1,
  attempted_at = @attempt_at,
  schedule_at = @next_schedule_at
FROM locked_jobs
WHERE u_jobs.event_id = locked_jobs.event_id 
RETURNING u_jobs.topic, u_jobs.event_id