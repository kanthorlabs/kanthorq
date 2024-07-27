package kanthorq

type ConsumerRegistry struct {
	Name       string `json:"name" validate:"required,is_collection_name"`
	StreamName string `json:"stream_name" validate:"required,is_collection_name"`
	Topic      string `json:"topic" validate:"required,is_topic"`
	Cursor     string `json:"cursor"`
	AttemptMax int16  `json:"attempt_max"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}
