package puller

import (
	"context"

	"github.com/kanthorlabs/kanthorq/core"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
)

var _ Puller = (*retry)(nil)

func NewRetry(
	cm pgcm.ConnectionManager,
	stream *entities.StreamRegistry,
	consumer *entities.ConsumerRegistry,
	in PullerIn,
) Puller {
	return &retry{cm: cm, stream: stream, consumer: consumer, in: in}
}

type retry struct {
	cm       pgcm.ConnectionManager
	stream   *entities.StreamRegistry
	consumer *entities.ConsumerRegistry
	in       PullerIn
}

func (puller *retry) Do(ctx context.Context) (*PullerOut, error) {
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

func (puller *retry) convert(ctx context.Context, out *PullerOut) error {
	req := &core.TaskStateTransitionReq{
		Consumer:  puller.consumer,
		FromState: entities.StateRetryable,
		ToState:   entities.StateRunning,
		Size:      puller.in.Size,
	}
	res, err := core.DoWithCM(ctx, req, puller.cm)
	if err != nil {
		return err
	}

	out.Tasks = res.Tasks
	out.EventIds = res.EventIds
	return nil
}

func (puller *retry) fulfill(ctx context.Context, out *PullerOut) error {
	// no event to fulfill, return early
	if len(out.EventIds) == 0 {
		return nil
	}

	req := &core.StreamGetEventReq{
		Stream:   puller.stream,
		EventIds: out.EventIds,
	}
	res, err := core.DoWithCM(ctx, req, puller.cm)
	if err != nil {
		return err
	}
	out.Events = res.Events
	return nil
}
