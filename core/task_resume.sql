--->>> task_resume
UPDATE %s
SET state = @state_running::SMALLINT, finalized_at = 0
WHERE 
  event_id IN (%s) AND state IN (@state_discarded::SMALLINT, @state_cancelled::SMALLINT)
RETURNING event_id;
---<<< task_resume