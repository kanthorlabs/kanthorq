--->>> api_task_mark_running_as_completed
UPDATE %s
SET state = @completed_state::SMALLINT , finalized_at = @finalized_at::BIGINT
WHERE 
  event_id IN (%s) AND state = @running_state::SMALLINT
RETURNING event_id;
---<<< api_task_mark_running_as_completed