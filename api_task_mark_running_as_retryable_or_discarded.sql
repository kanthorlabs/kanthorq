--->>> api_task_mark_running_as_retryable_or_discarded
UPDATE %s
SET 
  state = CASE 
            WHEN attempt_count > @attempt_max THEN @discarded_state::SMALLINT
            ELSE @retryable_state::SMALLINT END,
  schedule_at = EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000 + ((attempt_count ^ 4) + (attempt_count ^ 4) * (RANDOM() * 0.2 - 0.1)) * 60 * 1000,
  attempt_count = attempt_count + 1
WHERE 
  event_id IN (%s) AND state = @running_state::SMALLINT
RETURNING event_id;
---<<< api_task_mark_running_as_retryable_or_discarded