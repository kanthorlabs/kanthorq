package entities

import (
	"time"

	"github.com/kanthorlabs/kanthorq/pkg/xid"
)

func EventId() string {
	return xid.New("event")
}

func EventIdFromTime(t time.Time) string {
	return xid.NewWithTime("event", t)
}
func NewEvent(subject string, body []byte) *Event {
	return &Event{
		Id:        EventId(),
		Subject:   subject,
		Body:      body,
		Metadata:  make(Metadata),
		CreatedAt: time.Now().UnixMilli(),
	}
}

type Event struct {
	Id      string `json:"id" validate:"required"`
	Subject string `json:"subject" validate:"required,is_subject"`
	Body    []byte `json:"body" validate:"required"`

	// Metadata carry additional information about the event.
	Metadata Metadata `json:"metadata" validate:"required"`
	// CreatedAt is when the event record was created.
	CreatedAt int64 `json:"created_at" validate:"required,gt=0"`
}
