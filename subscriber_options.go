package kanthorq

type SubscriberOptions struct {
	StreamName         string `json:"stream_name" yaml:"stream_name" validate:"required,is_collection_name"`
	ConsumerName       string `json:"consumer_name" yaml:"consumer_name" validate:"required,is_collection_name"`
	ConsumerSubject    string `json:"consumer_subject" yaml:"consumer_subject" validate:"required,is_subject"`
	ConsumerAttemptMax int16  `json:"consumer_attempt_max" yaml:"consumer_attempt_max" validate:"required,gte=0"`

	HandleTimeout int64 `json:"handle_timeout" yaml:"handle_timeout" validate:"required,gte=3000"`
}
