package publisher

type Options struct {
	Connection string `json:"connection" yaml:"connection" validate:"required"`
	StreamName string `json:"stream_name" yaml:"stream_name" validate:"required,is_collection_name"`
}
