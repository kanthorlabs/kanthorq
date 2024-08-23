package publisher

import (
	"context"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
	"github.com/kanthorlabs/kanthorq/pkg/xvalidator"
)

// NewPublisher creates a new publisher that uses the default stream
func New(options *Options) (Publisher, error) {
	if err := xvalidator.Validate.Struct(options); err != nil {
		return nil, err
	}
	cm, err := pgcm.New(options.Connection)
	if err != nil {
		return nil, err
	}
	return &primary{options: options, cm: cm}, nil
}

type Publisher interface {
	Start(ctx context.Context) (err error)
	Stop(ctx context.Context) (err error)
	Send(ctx context.Context, events ...*entities.Event) error
}
