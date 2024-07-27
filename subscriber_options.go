package kanthorq

type SubscriberOptions struct {
	StreamName         string `json:"stream_name" yaml:"stream_name" validate:"required,is_collection_name"`
	ConsumerName       string `json:"consumer_name" yaml:"consumer_name" validate:"required,is_collection_name"`
	ConsumerTopic      string `json:"consumer_topic" yaml:"consumer_topic" validate:"required,is_topic"`
	ConsumerAttemptMax int16  `json:"consumer_attempt_max" yaml:"consumer_attempt_max" validate:"required,gte=0"`
}
