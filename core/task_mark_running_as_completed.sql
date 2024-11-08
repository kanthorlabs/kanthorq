--->>> task_mark_running_as_completed
UPDATE %s
SET state = @state_completed::SMALLINT, finalized_at = @finalized_at
WHERE 
  event_id IN (%s) AND state = @state_running::SMALLINT
RETURNING event_id;
---<<< task_mark_running_as_completed