package kanthorq

type Stream struct {
	Tier    string
	Topic   string
	EventId string
}

type Consumer struct {
	Name   string
	Topic  string
	Cursor string
}

type ConsumerJob struct {
	Tier    string
	Topic   string
	EventId string
}
