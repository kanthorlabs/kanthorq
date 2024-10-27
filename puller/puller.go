package puller

import (
	"context"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
	"go.uber.org/zap"
)

type PullerFactory func(
	logger *zap.Logger,
	cm pgcm.ConnectionManager,
	stream *entities.StreamRegistry,
	consumer *entities.ConsumerRegistry,
	in PullerIn,
) Puller

type Puller interface {
	Do(ctx context.Context) (*PullerOut, error)
}

type PullerIn struct {
	// Size is how many events you want to pull at one batch
	Size int `validate:"required,gt=0,lte=500"`
	// WaitingTime is how long you want to wait before pulling again
	// if you didn't get enough events in current batch
	WaitingTime int64 `validate:"gte=1000"`
}

type PullerOut struct {
	Tasks    map[string]*entities.Task
	Events   []*entities.Event
	EventIds []string
}
