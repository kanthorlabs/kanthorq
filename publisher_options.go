package kanthorq

type PublisherOptions struct {
	StreamName string `json:"stream_name" yaml:"stream_name" validate:"required,is_collection_name"`
}
