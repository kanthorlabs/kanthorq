package kanthorq

import (
	"context"
	"log"
	"time"

	"github.com/kanthorlabs/kanthorq/subscriber"
)

func Sub(ctx context.Context, options *subscriber.Options, handler subscriber.Handler) error {
	clients := make([]subscriber.Subscriber, 0)

	if primary, err := subscriber.New(options); err == nil {
		clients = append(clients, primary)
	} else {
		return err
	}

	if retry, err := subscriber.NewRetry(options); err == nil {
		clients = append(clients, retry)
	} else {
		return err
	}

	// defer stop all clients
	defer func() {
		stopctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		for _, client := range clients {
			if err := client.Stop(stopctx); err != nil {
				log.Println(err)
			}
		}
	}()

	// start all clients
	startctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	for _, client := range clients {
		if err := client.Start(startctx); err != nil {
			return err
		}
	}

	rctx, stop := context.WithCancel(ctx)
	defer stop()
	errc := make(chan error, 1)

	for _, client := range clients {
		go func(c subscriber.Subscriber) {
			if err := c.Receive(rctx, handler); err != nil {
				stop()
				// if one of the clients return error, stop all clients
				errc <- err
			}
		}(client)
	}

	return <-errc
}
