package kanthorq

import (
	"context"
	"errors"
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
	puller   Puller
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
		StreamName:            sub.options.StreamName,
		ConsumerName:          sub.options.ConsumerName,
		ConsumerSubjectFilter: sub.options.ConsumerSubjectFilter,
		ConsumerAttemptMax:    sub.options.ConsumerAttemptMax,
	}
	res, err := Do(ctx, req, conn)
	if err != nil {
		return
	}

	sub.stream = res.StreamRegistry
	sub.consumer = res.ConsumerRegistry
	sub.puller = &PullerDefault{
		cm:       sub.cm,
		stream:   sub.stream,
		consumer: sub.consumer,
		in:       PullerInDefault,
	}
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

func (sub *subscriber) Receive(ctx context.Context, handler SubscriberHandler) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// The Pulling Workflow
			out, err := sub.puller.Do(ctx)
			if err != nil {
				return err
			}
			log.Println("received", len(out.Events), "events")
			if len(out.Events) == 0 {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(time.Millisecond * 300):
					// wait for a while
					log.Println("no events, wait for a while")
				}
				continue
			}

			// the events are already sorted ascending by event id
			// and we should respect the order of events by executing events in order
			for _, event := range out.Events {
				msg := &Message{
					Event:    event,
					Task:     out.Tasks[event.Id],
					cm:       sub.cm,
					consumer: sub.consumer,
				}

				if err = handler(ctx, event); err != nil {
					log.Println(err)
					if err := msg.Nack(ctx); err != nil {
						log.Println(err)
					}
					continue
				}

				if err := msg.Nack(ctx); err != nil {
					log.Println(err)
				}
			}
		}
	}
}
