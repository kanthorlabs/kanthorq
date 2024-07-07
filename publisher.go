package kanthorq

import (
	"context"
	"errors"
)

type Publisher interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Send(ctx context.Context, events ...*Event) error
}

func NewPublisher() (Publisher, error) {
	return nil, errors.New("unimplemented")
}
