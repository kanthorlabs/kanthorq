package kanthorq

import (
	"context"

	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
)

var _ Puller = (*PullerRetry)(nil)

type PullerRetry struct {
	cm       pgcm.ConnectionManager
	stream   *StreamRegistry
	consumer *ConsumerRegistry
	in       *PullerIn
}

func (puller *PullerRetry) Do(ctx context.Context) (*PullerOut, error) {
	out := &PullerOut{
		Tasks:    make(map[string]*Task),
		Events:   make([]*Event, 0),
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

func (puller *PullerRetry) convert(ctx context.Context, out *PullerOut) error {
	req := &TaskStateTransitionReq{
		Consumer:  puller.consumer,
		FromState: StateRetryable,
		ToState:   StateRunning,
		Size:      puller.in.Size,
	}
	res, err := DoWithCM(ctx, req, puller.cm)
	if err != nil {
		return err
	}

	out.Tasks = res.Tasks
	out.EventIds = res.EventIds
	return nil
}

func (puller *PullerRetry) fulfill(ctx context.Context, out *PullerOut) error {
	// no event to fulfill, return early
	if len(out.EventIds) == 0 {
		return nil
	}

	req := &StreamGetEventReq{
		Stream:   puller.stream,
		EventIds: out.EventIds,
	}
	res, err := DoWithCM(ctx, req, puller.cm)
	if err != nil {
		return err
	}
	out.Events = res.Events
	return nil
}
