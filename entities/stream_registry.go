package entities

import (
	"time"

	"github.com/kanthorlabs/kanthorq/pkg/xid"
)

func StreamId() string {
	return xid.New("stream")
}

func StreamIdFromTime(t time.Time) string {
	return xid.NewWithTime("stream", t)
}

type StreamRegistry struct {
	Id        string `json:"id" validate:"required"`
	Name      string `json:"name" validate:"required,is_collection_name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
