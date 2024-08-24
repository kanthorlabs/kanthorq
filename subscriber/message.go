package subscriber

import (
	"context"
	"sync"

	"github.com/kanthorlabs/kanthorq/core"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
)

type Message struct {
	Event *entities.Event
	Task  *entities.Task

	cm       pgcm.ConnectionManager
	consumer *entities.ConsumerRegistry

	mu   sync.Mutex
	done bool
}

// Ack is safe to call multiple times
func (msg *Message) Ack(ctx context.Context) error {
	msg.mu.Lock()
	defer msg.mu.Unlock()

	if msg.done {
		return nil
	}
	msg.done = true

	req := &core.TaskMarkRunningAsCompletedReq{
		Consumer: msg.consumer,
		Tasks:    []*entities.Task{msg.Task},
	}

	// @TODO: if res.Noop has value, should log it here
	_, err := core.DoWithCM(ctx, req, msg.cm)
	return err
}

// Nack is safe to call multiple times
func (msg *Message) Nack(ctx context.Context) error {
	msg.mu.Lock()
	defer msg.mu.Unlock()

	if msg.done {
		return nil
	}
	msg.done = true

	req := &core.TaskMarkRunningAsRetryableOrDiscardedReq{
		Consumer: msg.consumer,
		Tasks:    []*entities.Task{msg.Task},
	}

	// @TODO: if res.Noop has value, should log it here
	_, err := core.DoWithCM(ctx, req, msg.cm)
	return err
}
