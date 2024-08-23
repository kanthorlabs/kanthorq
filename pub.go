package kanthorq

import (
	"context"
	"time"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/publisher"
)

func Pub(ctx context.Context, options *publisher.Options, ch chan *entities.Event) error {
	pub, err := publisher.New(options)
	if err != nil {
		return err
	}

	startctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	if err := pub.Start(startctx); err != nil {
		return err
	}
	defer func() {
		stopctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		if err := pub.Stop(stopctx); err != nil {
			panic(err)
		}
	}()

	for event := range ch {
		if err := pub.Send(ctx, event); err != nil {
			return err
		}
	}

	return nil
}
