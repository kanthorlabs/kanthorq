package subscriber

import "github.com/kanthorlabs/kanthorq/puller"

type Options struct {
	Connection                string          `json:"connection" yaml:"connection" validate:"required"`
	StreamName                string          `json:"stream_name" yaml:"stream_name" validate:"required,is_collection_name"`
	ConsumerName              string          `json:"consumer_name" yaml:"consumer_name" validate:"required,is_collection_name"`
	ConsumerSubjectIncludes   []string        `json:"consumer_subject_includes" yaml:"consumer_subject_includes" validate:"required,gt=0,dive,is_subject_filter"`
	ConsumerSubjectExcludes   []string        `json:"consumer_subject_excludes" yaml:"consumer_subject_excludes" validate:"gte=0,dive,is_subject_filter"`
	ConsumerAttemptMax        int16           `json:"consumer_attempt_max" yaml:"consumer_attempt_max" validate:"gte=0"`
	ConsumerVisibilityTimeout int64           `json:"consumer_visibility_timeout" yaml:"consumer_visibility_timeout" validate:"required,gte=1000"`
	Puller                    puller.PullerIn `json:"puller" yaml:"puller" validate:"required"`
}
