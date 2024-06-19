-- >>> consumer_job_state_change
WITH locked_jobs AS (
  SELECT
    event_id
  FROM %s AS l_jobs
  WHERE
    l_jobs.state = @from_state::SMALLINT
    AND l_jobs.schedule_at < @attempt_at::BIGINT
  ORDER BY
    l_jobs.state ASC,
    l_jobs.schedule_at ASC
  LIMIT @size::INTEGER
  FOR UPDATE SKIP LOCKED
)
UPDATE %s AS u_jobs
SET
  state = @to_state::SMALLINT,
  attempt_count = attempt_count + 1,
  attempted_at = @attempt_at::BIGINT,
  schedule_at = @next_schedule_at::BIGINT
FROM locked_jobs
WHERE u_jobs.event_id = locked_jobs.event_id 
RETURNING u_jobs.event_id
-- <<< consumer_job_state_change
