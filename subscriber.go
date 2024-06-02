package kanthorq

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kanthorlabs/common/clock"
	"github.com/kanthorlabs/kanthorq/api"
	"github.com/kanthorlabs/kanthorq/entities"
)

var _ Subscriber = (*subscriber)(nil)

func Sub(ctx context.Context, pool *pgxpool.Pool, clock clock.Clock, streamName, consumerName, topic string) (Subscriber, error) {
	consumer, err := Consumer(ctx, pool, &entities.Consumer{
		StreamName: streamName,
		Name:       consumerName,
		Topic:      topic,
	})
	if err != nil {
		return nil, err
	}

	return &subscriber{pool: pool, clock: clock, consumer: consumer}, nil
}

type Subscriber interface {
	Receive(ctx context.Context, handler SubscriberHandler, options ...SubscribeOption) (chan map[string]error, chan error)
}

type SubscriberHandler func(ctx context.Context, events map[string]*entities.StreamEvent) map[string]error

type subscriber struct {
	pool  *pgxpool.Pool
	clock clock.Clock

	consumer *entities.Consumer
}

func (sub *subscriber) Receive(ctx context.Context, handler SubscriberHandler, options ...SubscribeOption) (chan map[string]error, chan error) {
	var opts = &SubscriberOptions{
		Size:              DefaultSubscriberSize,
		VisibilityTimeout: DefaultSubscriberVisibilityTimeout,
	}
	for _, config := range options {
		config(opts)
	}

	var reportc = make(chan map[string]error, 100)
	var errorc = make(chan error, 1)

	go sub.receive(ctx, handler, opts, reportc, errorc)

	return reportc, errorc
}

func (sub *subscriber) receive(
	ctx context.Context,
	handler SubscriberHandler,
	opts *SubscriberOptions,
	reportc chan map[string]error,
	errorc chan error,
) {
	for {
		select {
		case <-ctx.Done():
			errorc <- ctx.Err()
		default:
			// We assume that the subscriber handler needs to process events for a long time,
			// hence, we should not hold a transaction too long.
			// Therefore, we move events from StateAvailable to StateRunning first,
			// then push them to the SubscriberHandler to handle the business logic.
			// If something goes wrong, we need to clean it up later.

			// @TODO: handle timeout
			subctx := context.Background()

			// pull job transaction
			tx, err := sub.pool.Begin(ctx)
			if err != nil {
				errorc <- err
			}

			c, err := api.ConsumerPull(sub.consumer, opts.Size).Do(ctx, tx)
			if err != nil {
				errorc <- errors.Join(err, tx.Rollback(ctx))
			}

			// there is no more job in stream
			if c.NextCursor == "" {
				if err := tx.Commit(ctx); err != nil {
					errorc <- err
				}
				// @TODO: sleep or do something to avoid busy loop
				continue
			}

			j, err := api.ConsumerJobPull(sub.consumer, opts.Size, opts.VisibilityTimeout).Do(ctx, tx, sub.clock)
			if err != nil {
				errorc <- errors.Join(err, tx.Rollback(ctx))
			}
			if err := tx.Commit(ctx); err != nil {
				errorc <- err
			}

			if len(j.Events) == 0 {
				// @TODO: sleep or do something to avoid busy loop
				continue
			}

			events := make(map[string]*entities.StreamEvent, len(j.Events))
			for _, event := range j.Events {
				events[event.EventId] = event
			}

			reports := handler(subctx, events)

			// comfirm job transaction
			tx, err = sub.pool.Begin(subctx)
			if err != nil {
				errorc <- err
			}

			var completed []string
			var retryable []string
			for _, event := range j.Events {
				if err, exist := reports[event.EventId]; exist && err != nil {
					retryable = append(retryable, event.EventId)
					continue
				}

				completed = append(completed, event.EventId)
			}

			// no error reports, mark jobs as completed
			if len(completed) > 0 {
				command := api.ConsumerJobMarkComplete(sub.consumer, completed)
				// mark complete may try to update not running job (because their state were updated in other transaction)
				// so the returning helps use determine what event status is updated successfully
				// TODO: report what update is successful
				if _, err := command.Do(subctx, tx, sub.clock); err != nil {
					errorc <- errors.Join(err, tx.Rollback(ctx))
					return
				}
			}

			if len(retryable) > 0 {
				// error reports, mark jobs as retryable
				command := api.ConsumerJobMarkRetry(sub.consumer, retryable)
				// same as ConsumerJobMarkRetry, the returning helps use determine what event status is updated successfully
				// TODO: report what update is successful
				if _, err := command.Do(subctx, tx, sub.clock); err != nil {
					errorc <- errors.Join(err, tx.Rollback(ctx))
					return
				}
			}

			if err := tx.Commit(subctx); err != nil {
				errorc <- err
			}

			reportc <- reports
		}
	}
}
