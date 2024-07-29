package kanthorq

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
)

var _ Subscriber = (*subscriber)(nil)

type subscriber struct {
	options *SubscriberOptions
	mu      sync.Mutex

	cm       pgcm.ConnectionManager
	stream   *StreamRegistry
	consumer *ConsumerRegistry
	puller   SubscriberPuller
}

func (sub *subscriber) Start(ctx context.Context) (err error) {
	sub.mu.Lock()
	defer sub.mu.Unlock()

	if err = sub.cm.Start(ctx); err != nil {
		return
	}

	conn, err := sub.cm.Connection(ctx)
	if err != nil {
		return
	}
	defer conn.Close(ctx)

	req := &ConsumerRegisterReq{
		StreamName:         sub.options.StreamName,
		ConsumerName:       sub.options.ConsumerName,
		ConsumerTopic:      sub.options.ConsumerTopic,
		ConsumerAttemptMax: sub.options.ConsumerAttemptMax,
	}
	res, err := Do(ctx, req, conn.Raw())
	if err != nil {
		return
	}

	sub.stream = res.StreamRegistry
	sub.consumer = res.ConsumerRegistry
	sub.puller = &DefaultSubscriberPuller{cm: sub.cm, stream: sub.stream, consumer: sub.consumer}
	return nil
}

func (sub *subscriber) Stop(ctx context.Context) (err error) {
	sub.mu.Lock()
	defer sub.mu.Unlock()

	if cmerr := sub.cm.Start(ctx); cmerr != nil {
		err = errors.Join(err, cmerr)
	}

	sub.stream = nil
	sub.consumer = nil
	return
}

func (sub *subscriber) Receive(ctx context.Context, handler SubscriberHandler) (err error) {

	for {
		// every round, we will set a timeout for current handler
		hctx, cancel := context.WithTimeout(ctx, time.Millisecond*time.Duration(sub.options.HandleTimeout))
		defer cancel()

		select {
		case <-hctx.Done():
			err = errors.Join(err, hctx.Err())
			return
		default:
			found, err := sub.handle(hctx, handler)
			if err != nil {
				// @TODO: log the error here
			}

			// no task was found, should wait for a while before next round
			if found > 0 {
				select {
				case <-time.After(time.Millisecond * time.Duration(sub.options.HandleInterval)):
					// do nothing, just wait for next round
					// @TODO: log message about no task found here
				case <-hctx.Done():
					// if context got canceled, should stop both the loop and the delay
					return errors.Join(err, hctx.Err())
				}
			}
		}
	}
}

func (sub *subscriber) handle(ctx context.Context, handler SubscriberHandler) (count int, err error) {
	// The Pulling Workflow
	out, err := sub.puller.Pull(ctx)
	if err != nil {
		return 0, err
	}

	// The Updating Workflow
	// @TODO: implement task logging
	succeed := []*Task{}
	failure := []*Task{}
	// the events are already sorted ascending by event id
	// and we should respect the order of events by executing events in order
	for _, event := range out.Events {
		if err = handler(ctx, event); err != nil {
			failure = append(failure, out.Tasks[event.Id])
			continue
		}

		succeed = append(succeed, out.Tasks[event.Id])
	}

	// we should run both complete and fail actions before report the error
	if cerr := sub.complete(ctx, succeed); cerr != nil {
		err = errors.Join(err, cerr)
	}
	if ferr := sub.fail(ctx, succeed); ferr != nil {
		err = errors.Join(err, ferr)
	}

	return len(out.Events), err
}

func (sub *subscriber) complete(ctx context.Context, tasks []*Task) (err error) {
	if len(tasks) == 0 {
		return
	}

	req := &TaskMarkRunningAsCompletedReq{
		Consumer: sub.consumer,
		Tasks:    tasks,
	}
	res, err := DoWithCM(ctx, req, sub.cm)
	if err == nil {
		return
	}

	if len(res.Noop) > 0 {
		// @TODO: report that some tasks were not updated
		return
	}

	return
}

func (sub *subscriber) fail(ctx context.Context, tasks []*Task) (err error) {
	if len(tasks) == 0 {
		return
	}

	req := &TaskMarkRunningAsRetryableOrDiscardedReq{
		Consumer: sub.consumer,
		Tasks:    tasks,
	}
	res, err := DoWithCM(ctx, req, sub.cm)
	if err == nil {
		return
	}

	if len(res.Noop) > 0 {
		// @TODO: report that some tasks were not updated
		return
	}

	return
}
