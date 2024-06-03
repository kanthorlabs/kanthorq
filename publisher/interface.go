package publisher

import (
	"context"

	"github.com/kanthorlabs/kanthorq/entities"
)

type Publisher interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error

	Send(ctx context.Context, events []*entities.StreamEvent) error
}
