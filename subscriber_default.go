package kanthorq

import (
	"context"
	"errors"
	"fmt"
	"log"
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
	receiver Receiver
}

func (sub *subscriber) Start(ctx context.Context) (err error) {
	sub.mu.Lock()
	defer sub.mu.Unlock()

	if err = sub.cm.Start(ctx); err != nil {
		return
	}

	conn, err := sub.cm.Acquire(ctx)
	if err != nil {
		return err
	}
	defer sub.cm.Release(ctx, conn)

	req := &ConsumerRegisterReq{
		StreamName:         sub.options.StreamName,
		ConsumerName:       sub.options.ConsumerName,
		ConsumerSubject:    sub.options.ConsumerSubject,
		ConsumerAttemptMax: sub.options.ConsumerAttemptMax,
	}
	res, err := Do(ctx, req, conn)
	if err != nil {
		return
	}

	sub.stream = res.StreamRegistry
	sub.consumer = res.ConsumerRegistry
	sub.receiver = &ReceiverDefault{cm: sub.cm, stream: sub.stream, consumer: sub.consumer}
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
				log.Println(err)
			}
			fmt.Printf("handled %d events\n", found)
		}
	}
}

func (sub *subscriber) handle(ctx context.Context, handler SubscriberHandler) (count int, err error) {
	// The Pulling Workflow
	// @TODO: remove hardcode
	out, err := sub.receiver.Pull(ctx, &ReceiverPullReq{
		Size:            100,
		ScanIntervalMax: 3,
	})
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
	if ferr := sub.fail(ctx, failure); ferr != nil {
		err = errors.Join(err, ferr)
	}
	if cerr := sub.complete(ctx, succeed); cerr != nil {
		err = errors.Join(err, cerr)
	}

	return len(out.Events), err
}

func (sub *subscriber) complete(ctx context.Context, tasks []*Task) error {
	if len(tasks) == 0 {
		return nil
	}

	req := &TaskMarkRunningAsCompletedReq{
		Consumer: sub.consumer,
		Tasks:    tasks,
	}
	res, err := DoWithCM(ctx, req, sub.cm)
	if err != nil {
		return err
	}

	if len(res.Noop) > 0 {
		// @TODO: report that some tasks were not updated
	}

	return nil
}

func (sub *subscriber) fail(ctx context.Context, tasks []*Task) error {
	if len(tasks) == 0 {
		return nil
	}

	req := &TaskMarkRunningAsRetryableOrDiscardedReq{
		Consumer: sub.consumer,
		Tasks:    tasks,
	}
	res, err := DoWithCM(ctx, req, sub.cm)
	if err != nil {
		return err
	}

	if len(res.Noop) > 0 {
		// @TODO: report that some tasks were not updated
	}

	return nil
}
