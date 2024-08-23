package puller

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/core"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
)

var _ Puller = (*primary)(nil)

type primary struct {
	cm       pgcm.ConnectionManager
	stream   *entities.StreamRegistry
	consumer *entities.ConsumerRegistry
	in       *PullerIn
}

func (puller *primary) Do(ctx context.Context) (*PullerOut, error) {
	out := &PullerOut{
		Tasks:    make(map[string]*entities.Task),
		Events:   make([]*entities.Event, 0),
		EventIds: make([]string, 0),
	}

	if err := puller.convert(ctx, out); err != nil {
		return nil, err
	}

	if err := puller.fulfill(ctx, out); err != nil {
		return nil, err
	}

	return out, nil
}

func (puller *primary) convert(ctx context.Context, out *PullerOut) error {
	conn, err := puller.cm.Acquire(ctx)
	if err != nil {
		return err
	}
	defer puller.cm.Release(ctx, conn)

	// there is no auto-rollback on context cancellation.
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	lock, err := (&core.ConsumerLockReq{Name: puller.consumer.Name}).Do(ctx, tx)
	if err != nil {
		// unable to lock consumer because it was using by another puller
		if errors.Is(err, pgx.ErrNoRows) {
			return tx.Rollback(ctx)
		}
		return errors.Join(err, tx.Rollback(ctx))
	}

	// IMPORTANT: make sure you only use the consumer that was locked successfully
	// otherwise you cannot get latest consumer cursor
	scan, err := (&core.StreamScanReq{Stream: puller.stream, Consumer: lock.Consumer, Size: puller.in.Size, WaitingTime: puller.in.WaitingTime}).Do(ctx, tx)
	if err != nil {
		return errors.Join(err, tx.Rollback(ctx))
	}

	if scan.Cursor != "" && scan.Cursor != lock.Consumer.Cursor {
		unlock, err := (&core.ConsumerUnlockReq{Name: lock.Consumer.Name, Cursor: scan.Cursor}).Do(ctx, tx)
		if err != nil {
			return errors.Join(err, tx.Rollback(ctx))
		}

		// override consumer with the updated one
		puller.consumer = unlock.Consumer
	}

	// no event to convert, return early
	if len(scan.Ids) == 0 {
		return tx.Commit(ctx)
	}

	convert, err := (&core.TaskConvertReq{Consumer: lock.Consumer, EventIds: scan.Ids, InitialState: entities.StateRunning}).Do(ctx, tx)
	if err != nil {
		return errors.Join(err, tx.Rollback(ctx))
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	out.Tasks = convert.Tasks
	out.EventIds = convert.EventIds
	return nil
}

func (puller *primary) fulfill(ctx context.Context, out *PullerOut) error {
	// no event to fulfill, return early
	if len(out.EventIds) == 0 {
		return nil
	}

	req := &core.StreamGetEventReq{
		Stream:   puller.stream,
		EventIds: out.EventIds,
	}
	res, err := core.DoWithCM(ctx, req, puller.cm)
	if err != nil {
		return err
	}
	out.Events = res.Events
	return nil
}
