package subscriber

import (
	"context"

	"github.com/kanthorlabs/kanthorq/core"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
)

type Message struct {
	Event *entities.Event
	Task  *entities.Task

	cm       pgcm.ConnectionManager
	consumer *entities.ConsumerRegistry
}

func (msg *Message) Ack(ctx context.Context) error {
	req := &core.TaskMarkRunningAsCompletedReq{
		Consumer: msg.consumer,
		Tasks:    []*entities.Task{msg.Task},
	}

	// @TODO: if res.Noop has value, should log it here
	_, err := core.DoWithCM(ctx, req, msg.cm)
	return err
}

func (msg *Message) Nack(ctx context.Context) error {
	req := &core.TaskMarkRunningAsRetryableOrDiscardedReq{
		Consumer: msg.consumer,
		Tasks:    []*entities.Task{msg.Task},
	}

	// @TODO: if res.Noop has value, should log it here
	_, err := core.DoWithCM(ctx, req, msg.cm)
	return err
}
