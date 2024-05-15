package kanthorq

type StreamMessage struct {
	StreamTopic string
	StreamName  string
	MessageId   string
}

type Consumer struct {
	Name         string
	StreamTopic  string
	StreamCursor string
}

type ConsumerMessage struct {
	StreamTopic string
	StreamName  string
	MessageId   string
}
