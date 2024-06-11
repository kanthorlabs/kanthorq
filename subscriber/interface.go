package subscriber

import (
	"context"

	"github.com/kanthorlabs/kanthorq/entities"
)

type Subscriber interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error

	Pull(ctx context.Context, options ...SubscribeOption) (map[string]*entities.StreamEvent, error)

	Consume(ctx context.Context, handler SubscriberHandler, options ...SubscribeOption)
	Failurec() <-chan map[string]error
	Error() <-chan error
}

type SubscriberHandler func(ctx context.Context, events map[string]*entities.StreamEvent) map[string]error
