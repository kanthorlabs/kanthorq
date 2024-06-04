package subscriber

import (
	"context"

	"github.com/kanthorlabs/kanthorq/entities"
)

type Subscriber interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error

	Report() <-chan map[string]error
	Error() <-chan error
	Consume(ctx context.Context, handler SubscriberHandler, options ...Option)
}

type SubscriberHandler func(ctx context.Context, events map[string]*entities.StreamEvent) map[string]error
