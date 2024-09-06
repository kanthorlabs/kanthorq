package kanthorq

import (
	"context"
	"reflect"
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

	startctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if err := pub.Start(startctx); err != nil {
		panic(err)
	}

	return pub, func() {
		stopctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		name := reflect.TypeOf(pub).Name()
		if err := pub.Stop(stopctx); err != nil {
			logger.Error("publisher stop with error", zap.String("publisher", name), zap.Error(err))
		}

		logger.Info("publisher stopped", zap.String("publisher", name))
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
		stopctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		for _, client := range clients {
			name := reflect.TypeOf(client).Name()
			if err := client.Stop(stopctx); err != nil {
				logger.Error("subscriber stop with error", zap.String("subscriber", name), zap.Error(err))
				return
			}

			logger.Info("subscriber stopped", zap.String("subscriber", name))
		}
	}()

	// start all clients
	startctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	for _, client := range clients {
		if err := client.Start(startctx); err != nil {
			return err
		}
	}

	rctx, stop := context.WithCancel(ctx)
	defer stop()

	for _, client := range clients {
		go func(c subscriber.Subscriber) {
			name := reflect.TypeOf(client).Name()
			if err := c.Receive(rctx, handler); err != nil {
				logger.Error("subscriber receive with error", zap.String("subscriber", name), zap.Error(err))
			}

			stop()
		}(client)
	}

	<-rctx.Done()
	return ctx.Err()
}
