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

var CollectionStreamEventProps = []string{
	"event_id",
	"topic",
	"body",
	"metadata",
	"created_at",
}

type StreamEvent struct {
	EventId   string         `json:"event_id"`
	Topic     string         `json:"topic"`
	Body      []byte         `json:"body"`
	Metadata  map[string]any `json:"metadata"`
	CreatedAt int64          `json:"created_at"`
}

var DefaultEventBodySize = 64 * 1024
