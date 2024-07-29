package kanthorq

type SubscriberOptions struct {
	StreamName         string `json:"stream_name" yaml:"stream_name" validate:"required,is_collection_name"`
	ConsumerName       string `json:"consumer_name" yaml:"consumer_name" validate:"required,is_collection_name"`
	ConsumerTopic      string `json:"consumer_topic" yaml:"consumer_topic" validate:"required,is_topic"`
	ConsumerAttemptMax int16  `json:"consumer_attempt_max" yaml:"consumer_attempt_max" validate:"required,gte=0"`

	HandleInterval int64 `json:"handle_interval" yaml:"handle_interval" validate:"required,gte=1000"`
	HandleTimeout  int64 `json:"handle_timeout" yaml:"handle_timeout" validate:"required,gte=3000"`
}
