package core

import "fmt"

var CollectionStream = "kanthorq_stream"
var CollectionConsumer = "kanthorq_consumer"

type Stream struct {
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func (ent *Stream) Properties() []string {
	return []string{"name", "created_at", "updated_at"}
}

type StreamEvent struct {
	Topic     string `json:"topic"`
	EventId   string `json:"event_id"`
	CreatedAt int64  `json:"created_at"`
}

func (ent *StreamEvent) Properties() []string {
	return []string{"topic", "event_id", "created_at"}
}

func StreamEventCollection(name string) string {
	return fmt.Sprintf("%s_%s", CollectionStream, name)
}

type Consumer struct {
	Name       string `json:"name"`
	StreamName string `json:"stream_name"`
	Topic      string `json:"topic"`
	Cursor     string `json:"cursor"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}

func (ent *Consumer) Properties() []string {
	return []string{"name", "stream_name", "topic", "cursor", "created_at", "updated_at"}
}

type ConsumerEvent struct {
	EventId   string `json:"event_id"`
	Name      string `json:"name"`
	Topic     string `json:"topic"`
	PullCount int16  `json:"pull_count"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func (ent *ConsumerEvent) Properties() []string {
	return []string{"event_id", "name", "topic", "pull_count", "created_at", "updated_at"}
}

func ConsumerEventCollection(name string) string {
	return fmt.Sprintf("%s_%s", CollectionConsumer, name)
}

type ConsumerCursor struct {
	Name    string `json:"name"`
	Current string `json:"current"`
	Next    string `json:"next"`
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
