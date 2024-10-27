--->>> task_mark_cancelled
UPDATE %s
SET state = @cancelled_state::SMALLINT, finalized_at = @finalized_at
WHERE 
  event_id IN (%s) AND state IN (@pending_state::SMALLINT, @available_state::SMALLINT, @retryable_state::SMALLINT)
RETURNING event_id;
---<<< task_mark_cancelled