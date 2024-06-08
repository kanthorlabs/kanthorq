package subscriber

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/api"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/q"
)

var _ Subscriber = (*subscriber)(nil)

func New(conf *Config) Subscriber {
	return &subscriber{
		conf: conf,
		// don't use unbuffer channel because it will block .Consume method
		// because when you send a value on an unbuffered channel,
		// the sending goroutine is blocked until another goroutine receives the value from the channel
		reportc: make(chan map[string]error, 1),
		errorc:  make(chan error, 1),
	}
}

type subscriber struct {
	conf    *Config
	reportc chan map[string]error
	errorc  chan error

	conn     *pgx.Conn
	consumer *entities.Consumer
}

func (sub *subscriber) Start(ctx context.Context) error {
	conn, err := pgx.Connect(ctx, sub.conf.ConnectionUri)
	if err != nil {
		return err
	}
	sub.conn = conn

	consumer, err := q.Consumer(ctx, sub.conn, &entities.Consumer{
		StreamName: sub.conf.StreamName,
		Topic:      sub.conf.Topic,
		Name:       sub.conf.ConsumerName,
	})
	if err != nil {
		return err
	}
	sub.consumer = consumer

	return nil
}

func (sub *subscriber) Stop(ctx context.Context) error {
	// wait for reportc and errorc to be closed
	<-sub.Error()
	<-sub.Report()
	return sub.conn.Close(ctx)
}

func (sub *subscriber) Error() <-chan error {
	return sub.errorc
}

func (sub *subscriber) Report() <-chan map[string]error {
	return sub.reportc
}

func (sub *subscriber) Consume(ctx context.Context, handler SubscriberHandler, options ...Option) {
	var opts = &Options{
		Size:              DefaultSize,
		VisibilityTimeout: DefaultVisibilityTimeout,
	}
	for _, configure := range options {
		configure(opts)
	}

	for {
		select {
		case <-ctx.Done():
			close(sub.errorc)
			close(sub.reportc)
			// @TODO: log context error
			return
			// @TODO: move channel handling to another place
		case err := <-sub.errorc:
			// @TODO: better error handler
			fmt.Println(err)
			continue
		case report := <-sub.reportc:
			// @TODO: better event error handler
			for eventId, err := range report {
				fmt.Println(eventId, err)
			}
		default:
			if ctx.Err() != nil {
				close(sub.errorc)
				close(sub.reportc)
				return
			}
			// We assume that the subscriber handler needs to process events for a long time,
			// hence, we should not hold a transaction too long.
			// Therefore, we move events from StateAvailable to StateRunning first,
			// then push them to the SubscriberHandler to handle the business logic.
			// If something goes wrong, we need to clean it up later.

			// @TODO: handle timeout
			subctx := context.Background()

			// pull job transaction
			tx, err := sub.conn.Begin(subctx)
			if err != nil {
				sub.errorc <- err
				continue
			}

			c, err := api.ConsumerPull(sub.consumer, opts.Size).Do(subctx, tx)
			if err != nil {
				sub.errorc <- errors.Join(err, tx.Rollback(subctx))
				continue
			}

			// there is no more job in stream
			if c.NextCursor == "" {
				if err := tx.Commit(subctx); err != nil {
					sub.errorc <- err
				}
				// @TODO: sleep or do something to avoid busy loop
				continue
			}

			j, err := api.ConsumerJobPull(sub.consumer, opts.Size, opts.VisibilityTimeout).Do(subctx, tx)
			if err != nil {
				sub.errorc <- errors.Join(err, tx.Rollback(subctx))
				continue
			}
			if err := tx.Commit(subctx); err != nil {
				sub.errorc <- err
				continue
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
			tx, err = sub.conn.Begin(subctx)
			if err != nil {
				sub.errorc <- err
				continue
			}

			var completed []string
			var retryable []string
			for _, event := range j.Events {
				if len(reports) == 0 {
					completed = append(completed, event.EventId)
					continue
				}

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
				if _, err := command.Do(subctx, tx); err != nil {
					sub.errorc <- errors.Join(err, tx.Rollback(ctx))
					continue
				}
			}

			if len(retryable) > 0 {
				// error reports, mark jobs as retryable
				command := api.ConsumerJobMarkRetry(sub.consumer, retryable)
				// same as ConsumerJobMarkRetry, the returning helps use determine what event status is updated successfully
				// TODO: report what update is successful
				if _, err := command.Do(subctx, tx); err != nil {
					sub.errorc <- errors.Join(err, tx.Rollback(ctx))
					continue
				}
			}

			if err := tx.Commit(subctx); err != nil {
				sub.errorc <- err
				continue
			}

			sub.reportc <- reports
		}
	}
}
