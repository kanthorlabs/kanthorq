package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/xlogger"
	"github.com/kanthorlabs/kanthorq/puller"
	"github.com/kanthorlabs/kanthorq/subscriber"
	"go.uber.org/zap"
)

func main() {
	var DATABASE_URI = "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable"
	if uri := os.Getenv("KANTHORQ_POSTGRES_URI"); uri != "" {
		DATABASE_URI = uri
	}

	var options = &subscriber.Options{
		// replace connection string with your database URI
		Connection: DATABASE_URI,
		// we use default stream for demo
		StreamName: entities.DefaultStreamName,
		// we use default consumer for demo
		ConsumerName: entities.DefaultConsumerName,
		// we will only receive events that match with the filter
		// so both system.say_hello and system.say_goodbye will be processed
		ConsumerSubjectIncludes: []string{"system.>"},
		// if task is failed, it will be retried it with this number of times
		ConsumerAttemptMax: entities.DefaultConsumerAttemptMax,
		// if task is stuck, we will wait this amount of time to reprocess it
		ConsumerVisibilityTimeout: entities.DefaultConsumerVisibilityTimeout,
		Puller: puller.PullerIn{
			// Size is how many events you want to pull at one batch
			Size: 100,
			// WaitingTime is how long you want to wait before finish current batch
			// because you don't get enough events defined in the Size attribute
			WaitingTime: 1000,
		},
	}

	// listen for SIGTERM so if you press Ctrl-C you can stop the program
	ctx, stop := signal.NotifyContext(context.TODO(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger := xlogger.New()
	// replace subscriber.New with different intialization method you get different subscriber
	// - subscriber.NewRetry for Retry Subscriber
	// - subscriber.NewAvailability for Availability Subscriber
	sub, err := subscriber.New(options, logger)
	if err != nil {
		panic(err)
	}

	var timeout = time.Second * 3

	// starting a subscriber should be use with timeout
	startctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if err := sub.Start(startctx); err != nil {
		panic(err)
	}

	defer func() {
		// graceful shutdown starting
		// don't reuse ctx here because it already done
		// you also need timeout here
		stopCtx, stopCancel := context.WithTimeout(context.Background(), timeout)
		defer stopCancel()
		if err := sub.Stop(stopCtx); err != nil {
			logger.Error("subscriber stop with error", zap.Error(err))
			return
		}
	}()

	// the main part, working on up comming events and tasks
	receiveCtx, receiveCancel := context.WithCancel(ctx)
	defer receiveCancel()

	// start receiving events and tasks
	go func() {
		err := sub.Receive(receiveCtx, func(ctx context.Context, msg *subscriber.Message) error {
			ts := time.UnixMilli(msg.Event.CreatedAt).Format(time.RFC3339)
			// print out recevied event
			fmt.Printf("RECEIVED: %s | %s | %s\n", msg.Event.Id, msg.Event.Subject, ts)
			return nil
		})

		if err != nil && !errors.Is(err, context.Canceled) {
			logger.Error("subscriber receive with error", zap.Error(err))
		}

		// subscriber is done, should cancel the context to trigger other workflows
		receiveCancel()
	}()

	<-receiveCtx.Done()
}
