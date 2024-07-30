package kanthorq

import (
	"strings"

	"github.com/kanthorlabs/kanthorq/pkg/idx"
)

func TaskId() string {
	return idx.New("task")
}

func TopicFilter(topic string) string {
	return strings.Replace(topic, "*", "%", -1)
}

type Task struct {
	EventId      string `json:"event_id" validate:"required"`
	Topic        string `json:"topic" validate:"required,is_topic"`
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
