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
	AttemptCount int16  `json:"attempt_count"`
	AttemptedAt  int64  `json:"attempted_at"`
	ScheduleAt   int64  `json:"schedule_at"`
	FinalizedAt  int64  `json:"finalized_at"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}

type JobState int

const (
	StateDiscarded JobState = -101
	StateCancelled JobState = -100
	StateAvailable JobState = 0
	StateRunning   JobState = 1
	StateCompleted JobState = 100
	StateRetryable JobState = 101
)

func (state JobState) String() string {
	switch state {
	case StateDiscarded:
		return "Discarded"
	case StateCancelled:
		return "Cancelled"
	case StateAvailable:
		return "Available"
	case StateRunning:
		return "Running"
	case StateCompleted:
		return "Completed"
	case StateRetryable:
		return "Retryable"
	default:
		return fmt.Sprintf("Unknown JobState (%d)", state)
	}
}
