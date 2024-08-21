--->>> api_task_mark_running_as_retryable_or_discarded
UPDATE %s
SET 
  state = CASE 
            WHEN attempt_count > @attempt_max THEN @discarded_state::SMALLINT
            ELSE @retryable_state::SMALLINT END,
  attempt_count = attempt_count + 1,
  schedule_at = EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000 + ((attempt_count ^ 4) * (1 + RANDOM() * 0.2 - 0.1)) * 1000
WHERE 
  event_id IN (%s) AND state = @running_state::SMALLINT
RETURNING event_id, state;
---<<< api_task_mark_running_as_retryable_or_discarded