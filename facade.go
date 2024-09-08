package kanthorq

import (
	"context"
	"errors"
	"time"

	"github.com/kanthorlabs/kanthorq/pkg/xlogger"
	"github.com/kanthorlabs/kanthorq/publisher"
	"github.com/kanthorlabs/kanthorq/subscriber"
	"go.uber.org/zap"
)

var timeout = time.Second * 3

func Pub(ctx context.Context, options *publisher.Options) (publisher.Publisher, func()) {
	logger := xlogger.New()

	pub, err := publisher.New(options, logger)
	if err != nil {
		panic(err)
	}

	startCtx, startCancel := context.WithTimeout(ctx, timeout)
	defer startCancel()
	if err := pub.Start(startCtx); err != nil {
		panic(err)
	}

	return pub, func() {
		stopCtx, stopCancel := context.WithTimeout(ctx, timeout)
		defer stopCancel()

		if err := pub.Stop(stopCtx); err != nil && !errors.Is(err, context.Canceled) {
			logger.Error("publisher stop with error", zap.Error(err))
		}

		logger.Info("publisher stopped")
	}
}

func Sub(ctx context.Context, options *subscriber.Options, handler subscriber.Handler) error {
	logger := xlogger.New()
	clients := make([]subscriber.Subscriber, 0)

	if primary, err := subscriber.New(options, logger); err == nil {
		clients = append(clients, primary)
	} else {
		return err
	}

	if retry, err := subscriber.NewRetry(options, logger); err == nil {
		clients = append(clients, retry)
	} else {
		return err
	}

	if visibility, err := subscriber.NewVisibility(options, logger); err == nil {
		clients = append(clients, visibility)
	} else {
		return err
	}

	//  stop all clients
	defer func() {
		// graceful shutdown starting
		// don't reuse ctx here because it already done
		// you also need timeout here
		stopCtx, stopCancel := context.WithTimeout(context.Background(), timeout)
		defer stopCancel()

		for _, client := range clients {
			if err := client.Stop(stopCtx); err != nil {
				logger.Error("subscriber stop with error", zap.Error(err))
				return
			}

		}
	}()

	// start all clients
	startCtx, startCancel := context.WithTimeout(ctx, timeout)
	defer startCancel()
	for _, client := range clients {
		if err := client.Start(startCtx); err != nil {
			return err
		}
	}

	receiveCtx, receiveCancel := context.WithCancel(ctx)
	defer receiveCancel()

	for _, client := range clients {
		go func(c subscriber.Subscriber) {
			if err := c.Receive(receiveCtx, handler); err != nil && !errors.Is(err, context.Canceled) {
				logger.Error("subscriber receive with error", zap.Error(err))
			}

			receiveCancel()
		}(client)
	}

	<-receiveCtx.Done()
	return ctx.Err()
}
