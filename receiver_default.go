package kanthorq

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
)

type ReceiverDefault struct {
	cm       pgcm.ConnectionManager
	stream   *StreamRegistry
	consumer *ConsumerRegistry
}

func (r *ReceiverDefault) Pull(ctx context.Context, req *ReceiverPullReq) (*ReceiverPullRes, error) {
	conn, err := r.cm.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer r.cm.Release(ctx, conn)

	// there is no auto-rollback on context cancellation.
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	lock, err := (&ConsumerLockReq{Name: r.consumer.Name}).Do(ctx, tx)
	if err != nil {
		return nil, errors.Join(err, tx.Rollback(ctx))
	}

	// IMPORTANT: make sure you only use the consumer that was locked successfully
	// otherwise you cannot get latest consumer cursor

	scan, err := (&StreamScanReq{Stream: r.stream, Consumer: lock.Consumer, Size: req.Size, IntervalMax: req.ScanIntervalMax}).Do(ctx, tx)
	if err != nil {
		return nil, errors.Join(err, tx.Rollback(ctx))
	}

	// TODO: scan has no id?

	convert, err := (&TaskConvertReq{Consumer: lock.Consumer, EventIds: scan.Ids, InitialState: StateRunning}).Do(ctx, tx)
	if err != nil {
		return nil, errors.Join(err, tx.Rollback(ctx))
	}

	// TODO: convert has no task?

	// must use event ids from convert step to make sure we only get events of successful converted tasks
	fulfill, err := (&StreamGetEventReq{Stream: r.stream, EventIds: convert.EventIds}).Do(ctx, tx)
	if err != nil {
		return nil, errors.Join(err, tx.Rollback(ctx))
	}

	// TODO: fullfill has no event?

	unlock, err := (&ConsumerUnlockReq{Name: lock.Consumer.Name, Cursor: scan.Cursor}).Do(ctx, tx)
	if err != nil {
		return nil, errors.Join(err, tx.Rollback(ctx))
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	// for tracking purposes only
	r.consumer = unlock.Consumer
	return &ReceiverPullRes{Events: fulfill.Events, Tasks: convert.Tasks}, nil
}
