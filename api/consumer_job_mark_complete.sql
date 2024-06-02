-- consumer_job_mark_complete
UPDATE %s
SET state = @complete_state, finalized_at = @finalized_at
WHERE 
  event_id IN (%s)
  -- make sure we only move job that are in running state to completed state
  AND state = @running_state
RETURNING event_id