package kanthorq

import (
	"reflect"
	"time"

	"github.com/kanthorlabs/kanthorq/pkg/idx"
)

func NewEvent(subject string, body []byte) *Event {
	return &Event{
		Id:        EventId(),
		Subject:   subject,
		Body:      body,
		Metadata:  make(Metadata),
		CreatedAt: time.Now().UnixMilli(),
	}
}

func EventId() string {
	return idx.New("event")
}

func EventIdFromTime(t time.Time) string {
	return idx.NewWithTime("event", t)
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

type Event struct {
	Id        string   `json:"id" validate:"required"`
	Subject   string   `json:"subject" validate:"required,is_subject"`
	Body      []byte   `json:"body" validate:"required"`
	Metadata  Metadata `json:"metadata" validate:"required"`
	CreatedAt int64    `json:"created_at"`
}
