package entities

import (
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

type Event struct {
	Id        string   `json:"id" validate:"required"`
	Subject   string   `json:"subject" validate:"required,is_subject"`
	Body      []byte   `json:"body" validate:"required"`
	Metadata  Metadata `json:"metadata" validate:"required"`
	CreatedAt int64    `json:"created_at"`
}
