package entities

var CollectionConsumer = "kanthorq_consumer"
var CollectionConsumerProps = []string{"name", "stream_name", "topic", "cursor", "created_at", "updated_at"}

type Consumer struct {
	Name       string `json:"name"`
	StreamName string `json:"stream_name"`
	Topic      string `json:"topic"`
	Cursor     string `json:"cursor"`
	AttemptMax int16  `json:"attempt_max"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}
