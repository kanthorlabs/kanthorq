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
	tasks, eventIds, err := r.convert(ctx, req)
	if err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return &ReceiverPullRes{Tasks: map[string]*Task{}, Events: []*Event{}}, nil
	}

	events, err := r.fulfill(ctx, eventIds)
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return &ReceiverPullRes{Tasks: map[string]*Task{}, Events: []*Event{}}, nil
	}

	// make sure we only return task has event
	if len(events) != len(tasks) {
		for _, event := range events {
			if _, ok := tasks[event.Id]; !ok {
				delete(tasks, event.Id)
			}
		}
	}

	// for tracking purposes only
	return &ReceiverPullRes{Tasks: tasks, Events: events}, nil
}

func (r *ReceiverDefault) convert(ctx context.Context, req *ReceiverPullReq) (map[string]*Task, []string, error) {
	conn, err := r.cm.Acquire(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer r.cm.Release(ctx, conn)

	// there is no auto-rollback on context cancellation.
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, nil, err
	}
	lock, err := (&ConsumerLockReq{Name: r.consumer.Name}).Do(ctx, tx)
	if err != nil {
		return nil, nil, errors.Join(err, tx.Rollback(ctx))
	}

	// IMPORTANT: make sure you only use the consumer that was locked successfully
	// otherwise you cannot get latest consumer cursor
	scan, err := (&StreamScanReq{Stream: r.stream, Consumer: lock.Consumer, Size: req.Size, IntervalMax: req.ScanIntervalMax}).Do(ctx, tx)
	if err != nil {
		return nil, nil, errors.Join(err, tx.Rollback(ctx))
	}

	if scan.Cursor != "" && scan.Cursor != lock.Consumer.Cursor {
		unlock, err := (&ConsumerUnlockReq{Name: lock.Consumer.Name, Cursor: scan.Cursor}).Do(ctx, tx)
		if err != nil {
			return nil, nil, errors.Join(err, tx.Rollback(ctx))
		}

		// override consumer with the updated one
		r.consumer = unlock.Consumer
	}

	// no event id, commit work then return
	if len(scan.Ids) == 0 {
		return map[string]*Task{}, make([]string, 0), tx.Commit(ctx)
	}

	convert, err := (&TaskConvertReq{Consumer: lock.Consumer, EventIds: scan.Ids, InitialState: StateRunning}).Do(ctx, tx)
	if err != nil {
		return nil, nil, errors.Join(err, tx.Rollback(ctx))
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, err
	}

	return convert.Tasks, convert.EventIds, nil
}

func (r *ReceiverDefault) fulfill(ctx context.Context, eventIds []string) ([]*Event, error) {
	fulfill, err := DoWithCM(ctx, &StreamGetEventReq{Stream: r.stream, EventIds: eventIds}, r.cm)
	if err != nil {
		return nil, err
	}
	if len(fulfill.Events) == 0 {
		return make([]*Event, 0), nil
	}

	return fulfill.Events, nil
}
