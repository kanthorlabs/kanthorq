package subscriber

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kanthorlabs/kanthorq/core"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
	"github.com/kanthorlabs/kanthorq/puller"
)

var _ Subscriber = (*primary)(nil)

type primary struct {
	mu sync.Mutex

	options *Options
	cm      pgcm.ConnectionManager
	pullerF puller.PullerFactory

	stream   *entities.StreamRegistry
	consumer *entities.ConsumerRegistry
	puller   puller.Puller
}

func (sub *primary) Start(ctx context.Context) (err error) {
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

	req := &core.ConsumerRegisterReq{
		StreamName:            sub.options.StreamName,
		ConsumerName:          sub.options.ConsumerName,
		ConsumerSubjectFilter: sub.options.ConsumerSubjectFilter,
		ConsumerAttemptMax:    sub.options.ConsumerAttemptMax,
	}
	res, err := core.Do(ctx, req, conn)
	if err != nil {
		return
	}

	sub.stream = res.StreamRegistry
	sub.consumer = res.ConsumerRegistry
	sub.puller = sub.pullerF(sub.cm, sub.stream, sub.consumer, sub.options.Puller)
	return nil
}

func (sub *primary) Stop(ctx context.Context) (err error) {
	sub.mu.Lock()
	defer sub.mu.Unlock()

	if cmerr := sub.cm.Stop(ctx); cmerr != nil {
		err = errors.Join(err, cmerr)
	}

	sub.stream = nil
	sub.consumer = nil
	sub.puller = nil
	return
}

func (sub *primary) Receive(ctx context.Context, handler Handler) error {
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

			var wg sync.WaitGroup
			for _, event := range out.Events {
				wg.Add(1)

				go func(hctx context.Context, e *entities.Event) {
					defer wg.Done()

					msg := &Message{
						Event:    event,
						Task:     out.Tasks[event.Id],
						cm:       sub.cm,
						consumer: sub.consumer,
					}

					if err = handler(hctx, event); err != nil {
						if nerr := msg.Nack(hctx); nerr != nil {
							log.Println(fmt.Errorf("failed to nack message: %w", errors.Join(err, nerr)))
						}
						return
					}

					if err := msg.Ack(hctx); err != nil {
						log.Println(fmt.Errorf("failed to ack message: %w", err))
					}
				}(ctx, event)
			}

			wg.Wait()
		}
	}
}
