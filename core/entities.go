package core

var CollectionStream = "kanthorq_stream"
var CollectionConsumer = "kanthorq_consumer"
var CollectionConsumerJob = "kanthorq_consumer_job"

type Stream struct {
	Tier    string
	Topic   string
	EventId string
}

func (ent *Stream) Properties() []string {
	return []string{"tier", "topic", "event_id"}
}

type Consumer struct {
	Name   string
	Topic  string
	Cursor string
}

func (ent *Consumer) Properties() []string {
	return []string{"name", "topic", "cursor"}
}

type ConsumerJob struct {
	Tier       string
	Topic      string
	EventId    string
	WriteCount int16
}

func (ent *Stream) ConsumerJob() []string {
	return []string{"tier", "topic", "event_id", "write_count"}
}
