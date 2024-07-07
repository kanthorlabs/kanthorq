package kanthorq

import (
	"context"
	"errors"
)

type Subscriber interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Send(ctx context.Context, events ...*Event) error
}

func NewSubscriber() (Subscriber, error) {
	return nil, errors.New("unimplemented")
}
