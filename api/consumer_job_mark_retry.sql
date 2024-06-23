-- >>> consumer_job_mark_retry
UPDATE %s
SET 
  state = CASE WHEN attempt_count > (SELECT attempt_max FROM kanthorq_consumer where name = @consumer_name)
              THEN @discarded_state::SMALLINT
              ELSE @retryable_state::SMALLINT END,
  schedule_at = EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000 + ((attempt_count ^ 4) + (attempt_count ^ 4) * (RANDOM() * 0.2 - 0.1)) * 60 * 1000,
  attempt_count = attempt_count + 1
WHERE 
  event_id IN (%s)
  -- make sure we only move job that are in running state to retryable state
  AND state = @running_state::SMALLINT
RETURNING event_id
-- <<< consumer_job_mark_retry