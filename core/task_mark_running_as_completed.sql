--->>> task_mark_running_as_completed
UPDATE %s
SET state = @completed_state::SMALLINT, finalized_at = @finalized_at
WHERE 
  event_id IN (%s) AND state = @running_state::SMALLINT
RETURNING event_id;
---<<< task_mark_running_as_completed