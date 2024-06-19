-- >>> consumer_job_mark_retry
UPDATE %s
SET 
  state = CASE WHEN attempt_count > @attempt_max::SMALLINT
              THEN @discarded_state::SMALLINT
              ELSE @retryable_state::SMALLINT END
WHERE 
  event_id IN (%s)
  -- make sure we only move job that are in running state to retryable state
  AND state = @running_state::SMALLINT
RETURNING event_id
-- <<< consumer_job_mark_retry