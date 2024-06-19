-- >>> consumer_job_mark_complete
UPDATE %s
SET state = @complete_state::SMALLINT , finalized_at = @finalized_at::BIGINT
WHERE 
  event_id IN (%s)
  -- make sure we only move job that are in running state to completed state
  AND state = @running_state::SMALLINT
RETURNING event_id
-- <<< consumer_job_mark_complete