package kanthorq

type SubscriberOptions struct {
	StreamName            string `json:"stream_name" yaml:"stream_name" validate:"required,is_collection_name"`
	ConsumerName          string `json:"consumer_name" yaml:"consumer_name" validate:"required,is_collection_name"`
	ConsumerSubjectFilter string `json:"consumer_subject_filter" yaml:"consumer_subject_filter" validate:"required,is_subject_filter"`
	ConsumerAttemptMax    int16  `json:"consumer_attempt_max" yaml:"consumer_attempt_max" validate:"required,gte=0"`

	HandlerTimeout int64 `json:"handler_timeout" yaml:"handler_timeout" validate:"required,gte=1000"`
}
