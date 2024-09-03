package subscriber

import (
	"context"
	"errors"
	"runtime/debug"
	"sync"
	"time"

	"github.com/kanthorlabs/kanthorq/core"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
)

type Message struct {
	Event *entities.Event
	Task  *entities.Task

	cm       pgcm.ConnectionManager
	consumer *entities.ConsumerRegistry

	mu     sync.Mutex
	acked  bool
	nacked bool
}

// Ack is safe to call multiple times
func (msg *Message) Ack(ctx context.Context) error {
	msg.mu.Lock()
	defer msg.mu.Unlock()

	if msg.nacked {
		return errors.New("message already nacked")
	}

	if msg.acked {
		return nil
	}
	msg.acked = true

	req := &core.TaskMarkRunningAsCompletedReq{
		Consumer: msg.consumer,
		Tasks:    []*entities.Task{msg.Task},
	}

	// @TODO: if res.Noop has value, should log it here
	_, err := core.DoWithCM(ctx, req, msg.cm)
	return err
}

// Nack is safe to call multiple times
func (msg *Message) Nack(ctx context.Context, reason error) error {
	msg.mu.Lock()
	defer msg.mu.Unlock()

	if msg.acked {
		return errors.New("message already acked")
	}

	if msg.nacked {
		return nil
	}
	msg.nacked = true

	req := &core.TaskMarkRunningAsRetryableOrDiscardedReq{
		Consumer: msg.consumer,
		Tasks:    []*entities.Task{msg.Task},
		Error: entities.AttemptedError{
			At:    time.Now().UnixMilli(),
			Error: reason.Error(),
			Stack: string(debug.Stack()),
		},
	}

	// @TODO: if res.Noop has value, should log it here
	if _, err := core.DoWithCM(ctx, req, msg.cm); err != nil {
		return err
	}
	return nil
}
