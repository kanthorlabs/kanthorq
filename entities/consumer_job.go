package entities

import "fmt"

func CollectionConsumerJob(name string) string {
	return fmt.Sprintf("%s_%s", CollectionConsumer, name)
}

var CollectionConsumerJobProps = []string{"event_id", "name", "topic", "pull_count", "created_at", "updated_at"}

type ConsumerJob struct {
	EventId   string `json:"event_id"`
	Name      string `json:"name"`
	Topic     string `json:"topic"`
	State     int16  `json:"state"`
	PullCount int16  `json:"pull_count"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type JobState int

const (
	Discarded JobState = -101
	Cancelled JobState = -100
	Available JobState = 0
	Running   JobState = 1
	Completed JobState = 100
	Retryable JobState = 101
)

func (js JobState) String() string {
	switch js {
	case Discarded:
		return "Discarded"
	case Cancelled:
		return "Cancelled"
	case Available:
		return "Available"
	case Running:
		return "Running"
	case Completed:
		return "Completed"
	case Retryable:
		return "Retryable"
	default:
		return fmt.Sprintf("Unknown JobState (%d)", js)
	}
}
