BEGIN;

CREATE TABLE IF NOT EXISTS kanthorq_job (
  org_id VARCHAR(128) NOT NULL,
  id VARCHAR(128) NOT NULL,
  created_at BIGINT NOT NULL,
  updated_at BIGINT NOT NULL DEFAULT 0,
  
  -- attempted_at is the time that the job was last worked
  attempted_at BIGINT NOT NULL DEFAULT 0,
  -- attempt_count is the attempt number of the job. Jobs are inserted at 0, the
  -- number is incremented to 1 the first time work its worked, and may
  -- increment further if it's either snoozed or errors.
  attempt_count SMALLINT NOT NULL DEFAULT 0,
  -- attempted_by is the set of client IDs that have worked this job.
  attempted_by VARCHAR(128) [],
  -- Errors is a set of errors that occurred when the job was worked, one for
  -- each attempt
  attempted_errors JSONB [],
	-- scheduled_at is when the job is scheduled to become available to be
	-- worked. Jobs default to running immediately, but may be scheduled
	-- for the future when they're inserted. They may also be scheduled for
	-- later because they were snoozed or because they errored and have
	-- additional retry attempts remaining.
  scheduled_at BIGINT NOT NULL,

  -- args represent the user's definition of a job, allowing them to store any relevant information.
  args JSONB NOT NULL,
  -- metadata is an attribute used by kanthorq internally to store data.
  metadata JSONB NOT NULL DEFAULT '{}' ::jsonb,
  -- state is described by an integer to facilitate easy maintenance and indicates the current status of the job.
  -- -101 - discarded
  -- -100 - cancelled
  -- 0 - available
  -- 1 - running
  -- 2 - retryable
  -- 100 - completed
  state SMALLINT NOT NULL DEFAULT 0,

  PRIMARY KEY (org_id, id)
);

COMMIT;