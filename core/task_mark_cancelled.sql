--->>> task_mark_cancelled
UPDATE %s
SET state = @state_cancelled::SMALLINT, finalized_at = @finalized_at
WHERE 
  event_id IN (%s) AND state IN (@state_pending::SMALLINT, @state_available::SMALLINT, @state_retryable::SMALLINT)
RETURNING event_id;
---<<< task_mark_cancelled