package kanthorq

import (
	"time"

	"github.com/kanthorlabs/kanthorq/pkg/idx"
)

var EventProps = []string{"id", "topic", "body", "metadata", "created_at"}

type Event struct {
	Id        string                 `json:"id"`
	Topic     string                 `json:"topic"`
	Body      []byte                 `json:"body"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt int64                  `json:"created_at"`
}

func NewEvent(topic string, body []byte) *Event {
	return &Event{
		Id:        idx.New("evt"),
		Topic:     topic,
		Body:      body,
		Metadata:  make(map[string]interface{}),
		CreatedAt: time.Now().UnixMilli(),
	}
}
