package kanthorq

import "github.com/kanthorlabs/kanthorq/pkg/idx"

func ConsumerId() string {
	return idx.New("consumer")
}

type ConsumerRegistry struct {
	StreamId   string `json:"stream_id" validate:"required"`
	StreamName string `json:"stream_name" validate:"required,is_collection_name"`
	Id         string `json:"id" validate:"required"`
	Name       string `json:"name" validate:"required,is_collection_name"`
	Topic      string `json:"topic" validate:"required,is_topic"`
	Cursor     string `json:"cursor"`
	AttemptMax int16  `json:"attempt_max"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}
