package puller

import (
	"context"

	"github.com/kanthorlabs/kanthorq/core"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
	"go.uber.org/zap"
)

var _ Puller = (*visibility)(nil)

func NewVisibility(
	logger *zap.Logger,
	cm pgcm.ConnectionManager,
	stream *entities.StreamRegistry,
	consumer *entities.ConsumerRegistry,
	in PullerIn,
) Puller {
	logger = logger.With(zap.String("puller", "visibility"))
	return &visibility{
		logger:   logger,
		cm:       cm,
		stream:   stream,
		consumer: consumer,
		in:       in,
	}
}

type visibility struct {
	logger   *zap.Logger
	cm       pgcm.ConnectionManager
	stream   *entities.StreamRegistry
	consumer *entities.ConsumerRegistry
	in       PullerIn
}

func (puller *visibility) Do(ctx context.Context) (*PullerOut, error) {
	out := &PullerOut{
		Tasks:    make(map[string]*entities.Task),
		Events:   make([]*entities.Event, 0),
		EventIds: make([]string, 0),
	}

	if err := puller.convert(ctx, out); err != nil {
		return nil, err
	}

	if err := puller.fulfill(ctx, out); err != nil {
		return nil, err
	}

	return out, nil
}

func (puller *visibility) convert(ctx context.Context, out *PullerOut) error {
	req := &core.TaskStateTransitionReq{
		Consumer:  puller.consumer,
		FromState: entities.StateRunning,
		ToState:   entities.StateRunning,
		Size:      puller.in.Size,
	}
	res, err := core.DoWithCM(ctx, puller.cm, req)
	if err != nil {
		return err
	}

	out.Tasks = res.Tasks
	out.EventIds = res.EventIds
	return nil
}

func (puller *visibility) fulfill(ctx context.Context, out *PullerOut) error {
	// no event to fulfill, return early
	if len(out.EventIds) == 0 {
		return nil
	}

	req := &core.StreamGetEventReq{
		Stream:   puller.stream,
		EventIds: out.EventIds,
	}
	res, err := core.DoWithCM(ctx, puller.cm, req)
	if err != nil {
		return err
	}
	out.Events = res.Events
	return nil
}
