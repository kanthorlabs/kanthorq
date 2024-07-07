package kanthorq

type Event struct {
	Id        string                 `json:"id"`
	Topic     string                 `json:"topic"`
	Body      []byte                 `json:"body"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt int64                  `json:"created_at"`
}
