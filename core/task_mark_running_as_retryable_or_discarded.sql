--->>> task_mark_running_as_retryable_or_discarded
UPDATE %s
SET 
  state = CASE 
    WHEN attempt_count >= @attempt_max THEN @discarded_state::SMALLINT
    ELSE @retryable_state::SMALLINT END,
  finalized_at = CASE
    WHEN attempt_count >= @attempt_max THEN @finalized_at::BIGINT
    ELSE finalized_at END
WHERE 
  event_id IN (%s) AND state = @running_state::SMALLINT
RETURNING event_id, state;
---<<< task_mark_running_as_retryable_or_discarded