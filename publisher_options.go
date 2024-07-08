package kanthorq

import "github.com/kanthorlabs/kanthorq/pkg/validator"

type PublisherOptions struct {
	StreamName string `json:"stream_name" yaml:"stream_name" validate:"required,is_collection_name"`
}

func (options *PublisherOptions) Validate() error {
	return validator.Validate.Struct(options)
}
