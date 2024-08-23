package publisher

type Options struct {
	StreamName string `json:"stream_name" yaml:"stream_name" validate:"required,is_collection_name"`
}
