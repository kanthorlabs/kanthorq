package kanthorq

import (
	"context"

	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
)

type DefaultSubscriberPuller struct {
	cm       pgcm.ConnectionManager
	stream   *StreamRegistry
	consumer *ConsumerRegistry
}

func (puller *DefaultSubscriberPuller) Pull(ctx context.Context) (*SubscriberPullerOut, error) {
	converting, err := puller.convert(ctx, 100)
	if err != nil {
		return nil, err
	}

	fulfilling, err := puller.fulfill(ctx, converting.EventIds)
	if err != nil {
		return nil, err
	}

	out := &SubscriberPullerOut{Tasks: converting.Tasks, Events: fulfilling.Events}
	return out, nil
}

func (puller *DefaultSubscriberPuller) convert(ctx context.Context, size int) (*TaskConvertFromEventRes, error) {
	req := &TaskConvertFromEventReq{
		Consumer:         puller.consumer,
		Size:             size,
		InitialTaskState: StateRunning,
	}
	conn, err := puller.cm.Connection(ctx)
	if err != nil {
		return nil, err
	}
	return Do(ctx, req, conn.Raw())
}

func (puller *DefaultSubscriberPuller) fulfill(ctx context.Context, ids []string) (*EventGetRes, error) {
	req := &EventGetReq{
		Stream:   puller.stream,
		EventIds: ids,
	}
	conn, err := puller.cm.Connection(ctx)
	if err != nil {
		return nil, err
	}
	return Do(ctx, req, conn.Raw())
}
