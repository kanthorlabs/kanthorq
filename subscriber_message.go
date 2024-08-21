package kanthorq

import (
	"context"

	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
)

type Message struct {
	Event *Event
	Task  *Task

	cm       pgcm.ConnectionManager
	consumer *ConsumerRegistry
}

func (msg *Message) Ack(ctx context.Context) error {
	req := &TaskMarkRunningAsCompletedReq{
		Consumer: msg.consumer,
		Tasks:    []*Task{msg.Task},
	}

	// @TODO: if res.Noop has value, should log it here
	_, err := DoWithCM(ctx, req, msg.cm)
	return err
}

func (msg *Message) Nack(ctx context.Context) error {
	req := &TaskMarkRunningAsRetryableOrDiscardedReq{
		Consumer: msg.consumer,
		Tasks:    []*Task{msg.Task},
	}

	// @TODO: if res.Noop has value, should log it here
	_, err := DoWithCM(ctx, req, msg.cm)
	return err
}
