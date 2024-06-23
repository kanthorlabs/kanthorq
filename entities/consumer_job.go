package entities

import "fmt"

func CollectionConsumerJob(name string) string {
	return fmt.Sprintf("%s_%s", CollectionConsumer, name)
}

var CollectionConsumerJobProps = []string{
	"event_id",
	"topic",
	"state",
	"attempt_count",
	"attempted_at",
	"schedule_at",
	"finalized_at",
	"created_at",
	"updated_at",
}

type ConsumerJob struct {
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

type JobState int

const (
	// JobStateDiscarded is the state for jobs that have errored enough times
	// that they're no longer eligible to be retried. Manual user invention
	// is required for them to be tried again.
	StateDiscarded JobState = -101
	// JobStateCancelled is the state for jobs that have been manually cancelled
	// by user request.
	StateCancelled JobState = -100
	// JobStateAvailable is the state for jobs that are immediately eligible to
	// be worked.
	StateAvailable JobState = 0
	// JobStateRunning are jobs which are actively running.
	StateRunning JobState = 1
	// JobStateScheduled is the state for jobs that are scheduled for the
	// future.
	StateCompleted JobState = 100
	// JobStateRetryable is the state for jobs that have errored, but will be
	// retried.
	StateRetryable JobState = 101
)

func (state JobState) String() string {
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
		return fmt.Sprintf("Unknown JobState (%d)", state)
	}
}
