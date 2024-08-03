package kanthorq

import (
	"context"

	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
	"github.com/kanthorlabs/kanthorq/pkg/validator"
)

func NewSubscriber(uri string, options *SubscriberOptions) (Subscriber, error) {
	if err := validator.Validate.Struct(options); err != nil {
		return nil, err
	}
	cm, err := pgcm.New(uri)
	if err != nil {
		return nil, err
	}
	return &subscriber{options: options, cm: cm}, nil
}

type Subscriber interface {
	Start(ctx context.Context) (err error)
	Stop(ctx context.Context) (err error)
	Receive(ctx context.Context, handler SubscriberHandler) error
}

type SubscriberHandler func(ctx context.Context, event *Event) error
