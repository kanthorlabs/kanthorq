package core

import "fmt"

var CollectionStream = "kanthorq_stream"
var CollectionConsumer = "kanthorq_consumer"

type Stream struct {
	Topic   string
	EventId string
}

func (ent *Stream) Properties() []string {
	return []string{"topic", "event_id"}
}

type Consumer struct {
	Name   string
	Topic  string
	Cursor string
}

func (ent *Consumer) Properties() []string {
	return []string{"name", "topic", "cursor"}
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
