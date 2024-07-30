package kanthorq

import "context"

type Receiver interface {
	Pull(ctx context.Context, req *ReceiverPullReq) (*ReceiverPullRes, error)
}

type ReceiverPullReq struct {
	MinSize        int
	ScanWindow     int64
	ScanRoundMax   int
	ScanRoundDelay int64
}

type ReceiverPullRes struct {
	Tasks  map[string]*Task
	Events []*Event
}
