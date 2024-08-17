package kanthorq

import "context"

type Puller interface {
	Do(ctx context.Context) (*PullerOut, error)
}

type PullerIn struct {
	Size        int   `validate:"required,gt=0"`
	WaitingTime int64 `validate:"gte=1000"`
}

type PullerOut struct {
	Tasks    map[string]*Task
	Events   []*Event
	EventIds []string
}

var PullerInDefault = &PullerIn{
	Size:        100,
	WaitingTime: 10000,
}
