package subscriber

import "github.com/kanthorlabs/kanthorq/puller"

type Options struct {
	StreamName            string           `json:"stream_name" yaml:"stream_name" validate:"required,is_collection_name"`
	ConsumerName          string           `json:"consumer_name" yaml:"consumer_name" validate:"required,is_collection_name"`
	ConsumerSubjectFilter []string         `json:"consumer_subject_filter" yaml:"consumer_subject_filter" validate:"required,gt=0,dive,is_subject_filter"`
	ConsumerAttemptMax    int16            `json:"consumer_attempt_max" yaml:"consumer_attempt_max" validate:"required,gte=0"`
	Puller                *puller.PullerIn `json:"puller" yaml:"puller" validate:"required"`
}
