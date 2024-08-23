package puller

import (
	"context"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
)

func New(
	cm pgcm.ConnectionManager,
	stream *entities.StreamRegistry,
	consumer *entities.ConsumerRegistry,
	in *PullerIn,
) Puller {
	return &primary{cm: cm, stream: stream, consumer: consumer, in: in}
}

func NewRetry(
	cm pgcm.ConnectionManager,
	stream *entities.StreamRegistry,
	consumer *entities.ConsumerRegistry,
	in *PullerIn,
) Puller {
	return &retry{cm: cm, stream: stream, consumer: consumer, in: in}
}

type Puller interface {
	Do(ctx context.Context) (*PullerOut, error)
}

type PullerIn struct {
	Size        int   `validate:"required,gt=0"`
	WaitingTime int64 `validate:"gte=1000"`
}

type PullerOut struct {
	Tasks    map[string]*entities.Task
	Events   []*entities.Event
	EventIds []string
}
