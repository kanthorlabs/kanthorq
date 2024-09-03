package publisher

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
)

type Publisher interface {
	Start(ctx context.Context) (err error)
	Stop(ctx context.Context) (err error)
	Send(ctx context.Context, events []*entities.Event) error
	SendTx(ctx context.Context, events []*entities.Event, tx pgx.Tx) error
}
