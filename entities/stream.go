package entities

import "fmt"

var CollectionStream = "kanthorq_stream"
var CollectionStreamProps = []string{"name", "created_at", "updated_at"}

type Stream struct {
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func CollectionStreamEvent(name string) string {
	return fmt.Sprintf("%s_%s", CollectionStream, name)
}

var CollectionStreamEventProps = []string{"topic", "event_id", "created_at"}

type StreamEvent struct {
	Topic     string `json:"topic"`
	EventId   string `json:"event_id"`
	CreatedAt int64  `json:"created_at"`
}

type EventPk struct {
	Topic   string `json:"topic"`
	EventId string `json:"event_id"`
}
