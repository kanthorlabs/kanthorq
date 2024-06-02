-- consumer_job_mark_retry
UPDATE %s
SET state = @retry_state
WHERE 
  event_id IN (%s)
  -- make sure we only move job that are in running state to completed state
  AND state = @running_state
RETURNING event_id