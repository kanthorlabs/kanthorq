package subscriber

import (
	"context"
	"errors"
	"runtime/debug"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/core"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
)

var (
	nacked = -1
	acked  = 1
)

type Message struct {
	Event *entities.Event
	Task  *entities.Task

	cm       pgcm.ConnectionManager
	consumer *entities.ConsumerRegistry

	mu           sync.Mutex
	acknowledged int
}

// Ack is safe to call multiple times
func (msg *Message) Ack(ctx context.Context) error {
	msg.mu.Lock()
	defer msg.mu.Unlock()

	acked, err := msg.acked()
	if err != nil {
		return err
	}
	if acked {
		return nil
	}

	req := &core.TaskMarkRunningAsCompletedReq{
		Consumer: msg.consumer,
		Tasks:    []*entities.Task{msg.Task},
	}

	_, err = core.DoWithCM(ctx, req, msg.cm)
	return err
}

func (msg *Message) AckTx(ctx context.Context, tx pgx.Tx) error {
	msg.mu.Lock()
	defer msg.mu.Unlock()

	acked, err := msg.acked()
	if err != nil {
		return err
	}
	if acked {
		return nil
	}

	req := &core.TaskMarkRunningAsCompletedReq{
		Consumer: msg.consumer,
		Tasks:    []*entities.Task{msg.Task},
	}

	_, err = req.Do(ctx, tx)
	return err
}

func (msg *Message) acked() (bool, error) {
	if msg.acknowledged == nacked {
		return false, errors.New("message is already nacked")
	}

	// already ack, don't do it again
	if msg.acknowledged != acked {
		msg.acknowledged = acked
		return false, nil
	}

	return true, nil
}

// Nack is safe to call multiple times
func (msg *Message) Nack(ctx context.Context, reason error) error {
	msg.mu.Lock()
	defer msg.mu.Unlock()

	nacked, err := msg.nacked()
	if err != nil {
		return err
	}
	if nacked {
		return nil
	}

	req := &core.TaskMarkRunningAsRetryableOrDiscardedReq{
		Consumer: msg.consumer,
		Tasks:    []*entities.Task{msg.Task},
		Error: entities.AttemptedError{
			At:    time.Now().UnixMilli(),
			Error: reason.Error(),
			Stack: string(debug.Stack()),
		},
	}

	_, err = core.DoWithCM(ctx, req, msg.cm)
	return err
}

func (msg *Message) NackTx(ctx context.Context, reason error, tx pgx.Tx) error {
	msg.mu.Lock()
	defer msg.mu.Unlock()

	nacked, err := msg.nacked()
	if err != nil {
		return err
	}
	if nacked {
		return nil
	}

	req := &core.TaskMarkRunningAsRetryableOrDiscardedReq{
		Consumer: msg.consumer,
		Tasks:    []*entities.Task{msg.Task},
		Error: entities.AttemptedError{
			At:    time.Now().UnixMilli(),
			Error: reason.Error(),
			Stack: string(debug.Stack()),
		},
	}

	_, err = req.Do(ctx, tx)
	return err
}

func (msg *Message) nacked() (bool, error) {
	if msg.acknowledged == acked {
		return false, errors.New("message is already acked")
	}

	if msg.acknowledged != nacked {
		msg.acknowledged = nacked
		return false, nil
	}

	return true, nil
}
