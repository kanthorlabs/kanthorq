package publisher

import (
	"context"

	"github.com/kanthorlabs/kanthorq/entities"
)

type Publisher interface {
	Send(ctx context.Context, events []*entities.StreamEvent) error
}
