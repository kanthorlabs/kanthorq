package subscriber

import (
	"context"

	"github.com/kanthorlabs/kanthorq/entities"
)

type Subscriber interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error

	Pull(ctx context.Context, options ...SubscribeOption) ([]*entities.StreamEvent, error)

	Consume(ctx context.Context, handler SubscriberHandler, options ...SubscribeOption)
	Failurec() <-chan map[string]error
	Error() <-chan error
}

type SubscriberHandler func(ctx context.Context, event *entities.StreamEvent) error
