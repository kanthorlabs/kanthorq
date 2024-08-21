--->>> api_task_mark_running_as_retryable_or_discarded
UPDATE %s
SET 
  state = CASE 
            WHEN attempt_count > @attempt_max THEN @discarded_state::SMALLINT
            ELSE @retryable_state::SMALLINT END,
  attempt_count = attempt_count + 1
WHERE 
  event_id IN (%s) AND state = @running_state::SMALLINT
RETURNING event_id, state;
---<<< api_task_mark_running_as_retryable_or_discarded