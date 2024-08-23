package subscriber

import (
	"context"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
	"github.com/kanthorlabs/kanthorq/pkg/validator"
)

func New(uri string, options *Options) (Subscriber, error) {
	if err := validator.Validate.Struct(options); err != nil {
		return nil, err
	}
	cm, err := pgcm.New(uri)
	if err != nil {
		return nil, err
	}
	return &primary{options: options, cm: cm}, nil
}

type Subscriber interface {
	Start(ctx context.Context) (err error)
	Stop(ctx context.Context) (err error)
	Receive(ctx context.Context, handler Handler) error
}

type Handler func(ctx context.Context, event *entities.Event) error
