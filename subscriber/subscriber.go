package subscriber

import (
	"context"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
	"github.com/kanthorlabs/kanthorq/pkg/xvalidator"
	"github.com/kanthorlabs/kanthorq/puller"
)

func New(options *Options) (Subscriber, error) {
	if err := xvalidator.Validate.Struct(options); err != nil {
		return nil, err
	}
	cm, err := pgcm.New(options.Connection)
	if err != nil {
		return nil, err
	}

	return &primary{options: options, cm: cm, pullerF: puller.New}, nil
}

func NewRetry(options *Options) (Subscriber, error) {
	if err := xvalidator.Validate.Struct(options); err != nil {
		return nil, err
	}
	cm, err := pgcm.New(options.Connection)
	if err != nil {
		return nil, err
	}

	return &primary{options: options, cm: cm, pullerF: puller.NewRetry}, nil
}

type Subscriber interface {
	Start(ctx context.Context) (err error)
	Stop(ctx context.Context) (err error)
	Receive(ctx context.Context, handler Handler) error
}

type Handler func(ctx context.Context, event *entities.Event) error
