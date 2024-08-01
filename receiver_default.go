package kanthorq

import (
	"context"

	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
)

type ReceiverDefault struct {
	cm       pgcm.ConnectionManager
	stream   *StreamRegistry
	consumer *ConsumerRegistry
}

func (r *ReceiverDefault) Pull(ctx context.Context, req *ReceiverPullReq) (*ReceiverPullRes, error) {
	converting, err := DoWithCM(ctx, &TaskConvertFromEventReq{
		Consumer:         r.consumer,
		InitialTaskState: StateRunning,
		Size:             req.Size,
		ScanWindow:       req.ScanWindow,
		ScanRoundMax:     req.ScanRoundMax,
		ScanRoundDelay:   req.ScanRoundDelay,
	}, r.cm)
	if err != nil {
		return nil, err
	}

	fulfilling, err := DoWithCM(ctx, &EventGetReq{
		Stream:   r.stream,
		EventIds: converting.EventIds,
	}, r.cm)
	if err != nil {
		return nil, err
	}

	res := &ReceiverPullRes{Tasks: converting.Tasks, Events: fulfilling.Events}
	return res, nil
}
