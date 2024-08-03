package kanthorq

import "context"

type Receiver interface {
	Pull(ctx context.Context, req *ReceiverPullReq) (*ReceiverPullRes, error)
}

type ReceiverPullReq struct {
	Size            int
	ScanIntervalMax int
}

type ReceiverPullRes struct {
	Tasks  map[string]*Task
	Events []*Event
}
