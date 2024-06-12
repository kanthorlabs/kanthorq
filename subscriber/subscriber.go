package subscriber

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

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
		failurec: make(chan map[string]error, 1),
		errorc:   make(chan error, 1),
	}
}

type subscriber struct {
	conf     *Config
	failurec chan map[string]error
	errorc   chan error
	mu       sync.Mutex

	conn     *pgx.Conn
	consumer *entities.Consumer
}

func (sub *subscriber) Start(ctx context.Context) error {
	if err := sub.connect(ctx); err != nil {
		return err
	}

	sub.mu.Lock()
	defer sub.mu.Unlock()
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

func (sub *subscriber) connect(ctx context.Context) error {
	sub.mu.Lock()
	defer sub.mu.Unlock()

	// @TODO: test what will happen if pgbouncer terminate a connection
	if sub.conn != nil && !sub.conn.IsClosed() {
		return nil
	}

	conn, err := pgx.Connect(ctx, sub.conf.ConnectionUri)
	if err != nil {
		return err
	}
	sub.conn = conn

	return nil
}

func (sub *subscriber) Stop(ctx context.Context) error {
	sub.mu.Lock()
	defer sub.mu.Unlock()

	// wait for failurec and errorc to be closed
	<-sub.Error()
	<-sub.Failurec()
	return sub.conn.Close(ctx)
}

func (sub *subscriber) Pull(ctx context.Context, options ...SubscribeOption) ([]*entities.StreamEvent, error) {
	var opts = &SubscribeOptions{
		Size:              DefaultSize,
		Timeout:           DefaultTimeout,
		VisibilityTimeout: DefaultVisibilityTimeout,
		WaitingTime:       DefaultWaitingTime,
	}
	for _, configure := range options {
		configure(opts)
	}

	ctx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()

	// both .Begin and .Rollback will teriminate the underlying connection
	// if the underlying connection is closed or context timeout
	tx, err := sub.conn.Begin(ctx)
	if err != nil {
		return nil, err
	}

	c, err := api.ConsumerPull(sub.consumer, opts.Size).Do(ctx, tx)
	// no new events to pull
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, tx.Rollback(ctx)
	}
	if err != nil {
		return nil, errors.Join(err, tx.Rollback(ctx))
	}

	// no more events to pull
	if c.NextCursor == "" {
		return nil, tx.Rollback(ctx)
	}

	j, err := api.ConsumerJobPull(sub.consumer, opts.Size, opts.VisibilityTimeout).Do(ctx, tx)
	// no new events to pull
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, tx.Rollback(ctx)
	}
	if err != nil {
		return nil, errors.Join(err, tx.Rollback(ctx))
	}

	// no event was found
	if len(j.Events) == 0 {
		return nil, tx.Rollback(ctx)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return j.Events, nil
}

func (sub *subscriber) Consume(ctx context.Context, handler SubscriberHandler, options ...SubscribeOption) {
	var opts = &SubscribeOptions{
		Size:              DefaultSize,
		Timeout:           DefaultTimeout,
		VisibilityTimeout: DefaultVisibilityTimeout,
		WaitingTime:       DefaultWaitingTime,
	}
	for _, configure := range options {
		configure(opts)
	}

	// start error handler
	go sub.errorh()
	// start failure handler
	go sub.failureh()

	for {
		select {
		case <-ctx.Done():
			close(sub.errorc)
			close(sub.failurec)
			return
		default:
			if ctx.Err() != nil {
				close(sub.errorc)
				close(sub.failurec)
				return
			}

			// both .Begin and .Rollback will teriminate the underlying connection
			// if the underlying connection is closed or context timeout
			// so we need an helper to check our connection status before start consuming
			if err := sub.connect(ctx); err != nil {
				close(sub.errorc)
				close(sub.failurec)
				// if we still can't connect, throw the error
				panic(err)
			}

			if err := sub.consume(handler, opts); err != nil {
				fmt.Printf("waiting for %s before retrying\n", opts.WaitingTime)
				time.Sleep(opts.WaitingTime)
			}
		}
	}
}

func (sub *subscriber) consume(handler SubscriberHandler, opts *SubscribeOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	events, err := sub.Pull(ctx)
	if err != nil {
		return sub.error(err)
	}

	failures := handler(ctx, events)

	// both .Begin and .Rollback will teriminate the underlying connection
	// if the underlying connection is closed or context timeout
	tx, err := sub.conn.Begin(ctx)
	if err != nil {
		return sub.error(err)
	}

	var completed []string
	var retryable []string
	for _, event := range events {
		if len(failures) == 0 {
			completed = append(completed, event.EventId)
			continue
		}

		if err, exist := failures[event.EventId]; exist && err != nil {
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
		if _, err := command.Do(ctx, tx); err != nil {
			return sub.error(err, tx.Rollback(ctx))
		}
	}

	if len(retryable) > 0 {
		// error reports, mark jobs as retryable
		command := api.ConsumerJobMarkRetry(sub.consumer, retryable)
		// same as ConsumerJobMarkRetry, the returning helps use determine what event status is updated successfully
		// TODO: report what update is successful
		if _, err := command.Do(ctx, tx); err != nil {
			return sub.error(err, tx.Rollback(ctx))
		}
	}

	if err := sub.error(tx.Commit(ctx)); err != nil {
		return err
	}
	sub.failurec <- failures
	return nil
}

func (sub *subscriber) Failurec() <-chan map[string]error {
	return sub.failurec
}

func (sub *subscriber) Error() <-chan error {
	return sub.errorc
}

func (sub *subscriber) error(errs ...error) error {
	var merged error

	for _, err := range errs {
		if err == nil {
			continue
		}
		merged = errors.Join(merged, err)
		sub.errorc <- err
	}

	return merged
}

func (sub *subscriber) errorh() {
	for err := range sub.errorc {
		fmt.Println("===errorh", err)
	}
}

func (sub *subscriber) failureh() {
	for failures := range sub.failurec {
		for eventId, err := range failures {
			fmt.Println("failureh===", eventId, err)
		}
	}
}
