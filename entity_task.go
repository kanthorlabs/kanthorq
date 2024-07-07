package kanthorq

import "fmt"

type Task struct {
	EventId      string `json:"event_id"`
	Topic        string `json:"topic"`
	State        int16  `json:"state"`
	ScheduleAt   int64  `json:"schedule_at"`
	AttemptCount int16  `json:"attempt_count"`
	AttemptedAt  int64  `json:"attempted_at"`
	FinalizedAt  int64  `json:"finalized_at"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}

type TaskState int

const (
	// StateDiscarded is the state for jobs that have errored enough times
	// that they're no longer eligible to be retried. Manual user invention
	// is required for them to be tried again.
	StateDiscarded TaskState = -102
	// StateCancelled is the state for jobs that have been manually cancelled
	// by user request.
	StateCancelled TaskState = -101
	// StateAvailable is the state for jobs that are immediately eligible to
	// be worked.
	StateAvailable TaskState = 1
	// StateRunning is the state for jobs jobs which are actively running.
	StateRunning TaskState = 2
	// Completed is the state for jobs that have successfully run to
	// completion.
	StateCompleted TaskState = 101
	// StateRetryable is the state for jobs that have errored, but will be
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
		return fmt.Sprintf("Unknown State (%d)", state)
	}
}
