--->>> consumer_job_mark_cancelled
UPDATE %s
SET state = @cancelled_state::SMALLINT , finalized_at = @finalized_at::BIGINT
WHERE 
  event_id IN (%s)
  -- make sure we only move jobs that are in (avaiable|running|retryable) state to completed state
  AND state in (@avaiable_state::SMALLINT, @running_state::SMALLINT, @retryable_state::SMALLINT)
RETURNING event_id
---<<< consumer_job_mark_cancelled