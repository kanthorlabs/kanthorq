package subscriber

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/kanthorlabs/kanthorq/core"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
	"github.com/kanthorlabs/kanthorq/puller"
	"go.uber.org/zap"
)

var _ Subscriber = (*primary)(nil)

type primary struct {
	mu sync.Mutex

	options *Options
	logger  *zap.Logger
	pullerF puller.PullerFactory
	cm      pgcm.ConnectionManager

	stream   *entities.StreamRegistry
	consumer *entities.ConsumerRegistry
	puller   puller.Puller
}

func (sub *primary) Start(ctx context.Context) error {
	sub.mu.Lock()
	defer sub.mu.Unlock()

	cm, err := pgcm.New(sub.options.Connection)
	if err != nil {
		return err
	}
	sub.cm = cm
	if err = sub.cm.Start(ctx); err != nil {
		return err
	}

	conn, err := sub.cm.Acquire(ctx)
	if err != nil {
		return err
	}
	defer sub.cm.Release(ctx, conn)

	req := &core.ConsumerRegisterReq{
		StreamName:                sub.options.StreamName,
		ConsumerName:              sub.options.ConsumerName,
		ConsumerSubjectIncludes:   sub.options.ConsumerSubjectIncludes,
		ConsumerSubjectExcludes:   sub.options.ConsumerSubjectExcludes,
		ConsumerAttemptMax:        sub.options.ConsumerAttemptMax,
		ConsumerVisibilityTimeout: sub.options.ConsumerVisibilityTimeout,
	}
	res, err := core.Do(ctx, req, conn)
	if err != nil {
		return err
	}

	sub.stream = res.StreamRegistry
	sub.consumer = res.ConsumerRegistry
	sub.logger.Info("started")

	sub.puller = sub.pullerF(sub.logger, sub.cm, sub.stream, sub.consumer, sub.options.Puller)
	return nil
}

func (sub *primary) Stop(ctx context.Context) error {
	sub.mu.Lock()
	defer sub.mu.Unlock()

	var err error

	if sub.cm != nil {
		if cmerr := sub.cm.Stop(ctx); cmerr != nil {
			err = errors.Join(err, cmerr)
		}
	}

	sub.stream = nil
	sub.consumer = nil
	sub.puller = nil
	sub.logger.Info("stopped")
	return err
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
			if len(out.Events) == 0 {
				continue
			}

			var wg sync.WaitGroup
			for _, event := range out.Events {
				wg.Add(1)

				msg := &Message{
					Event:    event,
					Task:     out.Tasks[event.Id],
					cm:       sub.cm,
					consumer: sub.consumer,
				}
				go sub.handle(ctx, handler, msg, &wg)
			}
			wg.Wait()
		}
	}
}

func (sub *primary) handle(ctx context.Context, handler Handler, msg *Message, wg *sync.WaitGroup) {
	defer sub.panic(ctx, msg, wg)

	if err := handler(ctx, msg); err != nil {
		if nerr := msg.Nack(ctx, err); nerr != nil {
			sub.logger.Error("failed to ack message", zap.Error(err))
		}
		return
	}

	if err := msg.Ack(ctx); err != nil {
		sub.logger.Error("failed to ack message", zap.Error(err))
	}
}

func (sub *primary) panic(ctx context.Context, msg *Message, wg *sync.WaitGroup) {
	if r := recover(); r != nil {
		var reason error
		if e, ok := r.(error); ok {
			reason = e
		} else {
			reason = fmt.Errorf("%v", r)
		}
		sub.logger.Error("catched panic", zap.Error(reason))

		if err := msg.Nack(ctx, reason); err != nil {
			sub.logger.Error("failed to nack message", zap.Error(err))
		}
	}

	wg.Done()
}
