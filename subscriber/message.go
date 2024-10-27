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

type msgstate = int16

var (
	nacked  msgstate = -1
	pending msgstate = 0
	acked   msgstate = 1
)

var (
	ErrAlreadyAcknowledged = errors.New("already acknowledged")
	ErrInvalidMessageState = errors.New("invalid message state")
)

type Message struct {
	Event *entities.Event
	Task  *entities.Task

	cm       pgcm.ConnectionManager
	consumer *entities.ConsumerRegistry

	mu           sync.Mutex
	acknowledged msgstate
}

// Ack is safe to call multiple times
func (msg *Message) Ack(ctx context.Context) error {
	msg.mu.Lock()
	defer msg.mu.Unlock()

	if err := msg.acked(); err != nil {
		if errors.Is(err, ErrAlreadyAcknowledged) {
			return nil
		}
		return err
	}

	req := &core.TaskMarkRunningAsCompletedReq{
		Consumer: msg.consumer,
		Tasks:    []*entities.Task{msg.Task},
	}

	_, err := core.DoWithCM(ctx, msg.cm, req)
	return err
}

func (msg *Message) AckTx(ctx context.Context, tx pgx.Tx) error {
	msg.mu.Lock()
	defer msg.mu.Unlock()

	if err := msg.acked(); err != nil {
		if errors.Is(err, ErrAlreadyAcknowledged) {
			return nil
		}
		return err
	}

	req := &core.TaskMarkRunningAsCompletedReq{
		Consumer: msg.consumer,
		Tasks:    []*entities.Task{msg.Task},
	}

	_, err := req.Do(ctx, tx)
	return err
}

func (msg *Message) acked() error {
	// safely call multiple times
	if msg.acknowledged == acked {
		return ErrAlreadyAcknowledged
	}

	// only allow move from pending to acked
	if msg.acknowledged == pending {
		msg.acknowledged = acked
		return nil
	}

	return ErrInvalidMessageState
}

// Nack is safe to call multiple times
func (msg *Message) Nack(ctx context.Context, reason error) error {
	msg.mu.Lock()
	defer msg.mu.Unlock()

	if err := msg.nacked(); err != nil {
		if errors.Is(err, ErrAlreadyAcknowledged) {
			return nil
		}
		return err
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

	_, err := core.DoWithCM(ctx, msg.cm, req)
	return err
}

func (msg *Message) NackTx(ctx context.Context, reason error, tx pgx.Tx) error {
	msg.mu.Lock()
	defer msg.mu.Unlock()

	if err := msg.nacked(); err != nil {
		if errors.Is(err, ErrAlreadyAcknowledged) {
			return nil
		}
		return err
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

	_, err := req.Do(ctx, tx)
	return err
}

func (msg *Message) nacked() error {
	// safely call multiple times
	if msg.acknowledged == nacked {
		return ErrAlreadyAcknowledged
	}

	// only allow move from pending to nnacked
	if msg.acknowledged == pending {
		msg.acknowledged = nacked
		return nil
	}

	return ErrInvalidMessageState
}
