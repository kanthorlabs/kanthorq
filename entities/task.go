package entities

import (
	"github.com/kanthorlabs/kanthorq/pkg/xid"
)

func TaskId() string {
	return xid.New("task")
}

type Task struct {
	EventId string `json:"event_id" validate:"required"`
	Subject string `json:"subject" validate:"required,is_subject"`

	// State is the state of task like `available` or `completed`.
	State int16 `json:"state"`
	// ScheduledAt is when the task is scheduled to become available to be
	// worked. Tasks default to running immediately, but may be scheduled
	// for the future when they're inserted. They may also be scheduled for
	// later because they were snoozed or because they errored and have
	// additional retry attempts remaining.
	ScheduleAt int64 `json:"schedule_at"`
	// AttemptCount is the attempt number of the task. Tasks are inserted at 0, the
	// number is incremented to 1 the first time work its worked, and may
	// increment further if it's either snoozed or errors.
	AttemptCount int16 `json:"attempt_count"`
	// AttemptedAt is the time that the task was last worked. Starts out as NOW()
	// on a new insert.
	AttemptedAt int64 `json:"attempted_at"`
	// AttemptedBy is the set of client IDs that have worked this task.
	AttemptedBy []string
	// FinalizedAt is the time at which the task was "finalized", meaning it was
	// either completed successfully or errored for the last time such that
	// it'll no longer be retried.
	FinalizedAt int64 `json:"finalized_at"`
	// CreatedAt is when the task record was created.
	CreatedAt int64 `json:"created_at"`
	// CreatedAt is when the task record was updated.
	UpdatedAt int64 `json:"updated_at"`
}

type TaskState int

const (
	// StateDiscarded is the state for tasks that have errored enough times
	// that they're no longer eligible to be retried. Manual user invention
	// is required for them to be tried again.
	StateDiscarded TaskState = -102
	// StateCancelled is the state for tasks that have been manually cancelled
	// by user request.
	StateCancelled TaskState = -101
	// StatePending is a state for tasks to be parked while waiting for some
	// external action before they can be worked. Tasks in pending will never be
	// worked or deleted unless moved out of this state by the user.
	StatePending TaskState = 0
	// StateAvailable is the state for tasks that are immediately eligible to
	// be worked.
	StateAvailable TaskState = 1
	// StateRunning is the state for tasks tasks which are actively running.
	StateRunning TaskState = 2
	// Completed is the state for tasks that have successfully run to
	// completion.
	StateCompleted TaskState = 101
	// StateRetryable is the state for tasks that have errored, but will be
	// retried.
	StateRetryable TaskState = 102
)

func (state TaskState) String() string {
	switch state {
	case StateDiscarded:
		return "discarded"
	case StateCancelled:
		return "cancelled"
	case StateAvailable:
		return "available"
	case StateRunning:
		return "running"
	case StateCompleted:
		return "completed"
	case StateRetryable:
		return "retryable"
	default:
		return ""
	}
}
