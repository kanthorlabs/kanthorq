-- >>> consumer_job_mark_retry
UPDATE %s
SET 
  SET state = CASE WHEN attempt_count > @attempt_max
              THEN @discarded_state
              ELSE @retryable_state END;
WHERE 
  event_id IN (%s)
  -- make sure we only move job that are in running state to retryable state
  AND state = @running_state
RETURNING event_id
-- <<< consumer_job_mark_retry