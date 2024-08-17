package kanthorq

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
)

var _ Puller = (*PullerDefault)(nil)

type PullerDefault struct {
	cm       pgcm.ConnectionManager
	stream   *StreamRegistry
	consumer *ConsumerRegistry
	in       *PullerIn
}

func (puller *PullerDefault) Do(ctx context.Context) (*PullerOut, error) {
	out := &PullerOut{
		Tasks:    make(map[string]*Task),
		Events:   make([]*Event, 0),
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

func (puller *PullerDefault) convert(ctx context.Context, out *PullerOut) error {
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
	lock, err := (&ConsumerLockReq{Name: puller.consumer.Name}).Do(ctx, tx)
	if err != nil {
		return errors.Join(err, tx.Rollback(ctx))
	}

	// IMPORTANT: make sure you only use the consumer that was locked successfully
	// otherwise you cannot get latest consumer cursor
	scan, err := (&StreamScanReq{Stream: puller.stream, Consumer: lock.Consumer, Size: puller.in.Size, WaitingTime: puller.in.WaitingTime}).Do(ctx, tx)
	if err != nil {
		return errors.Join(err, tx.Rollback(ctx))
	}

	if scan.Cursor != "" && scan.Cursor != lock.Consumer.Cursor {
		unlock, err := (&ConsumerUnlockReq{Name: lock.Consumer.Name, Cursor: scan.Cursor}).Do(ctx, tx)
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

	convert, err := (&TaskConvertReq{Consumer: lock.Consumer, EventIds: scan.Ids, InitialState: StateRunning}).Do(ctx, tx)
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

func (puller *PullerDefault) fulfill(ctx context.Context, out *PullerOut) error {
	// no event to fulfill, return early
	if len(out.EventIds) == 0 {
		return nil
	}

	fulfill, err := DoWithCM(ctx, &StreamGetEventReq{Stream: puller.stream, EventIds: out.EventIds}, puller.cm)
	if err != nil {
		return err
	}
	out.Events = fulfill.Events
	return nil
}
