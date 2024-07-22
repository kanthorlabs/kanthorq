package kanthorq

import (
	"reflect"
	"time"

	"github.com/kanthorlabs/kanthorq/pkg/idx"
)

func NewEvent(topic string, body []byte) *Event {
	return &Event{
		Id:        idx.New("evt"),
		Topic:     topic,
		Body:      body,
		Metadata:  make(map[string]interface{}),
		CreatedAt: time.Now().UnixMilli(),
	}
}

type Event struct {
	Id        string                 `json:"id"`
	Topic     string                 `json:"topic"`
	Body      []byte                 `json:"body"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt int64                  `json:"created_at"`
}

func EventProps() []string {
	var props []string
	eventType := reflect.TypeOf(Event{})

	for i := 0; i < eventType.NumField(); i++ {
		field := eventType.Field(i)
		props = append(props, field.Tag.Get("json"))
	}

	return props
}
