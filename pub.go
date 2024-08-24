package kanthorq

import (
	"context"
	"log"
	"time"

	"github.com/kanthorlabs/kanthorq/publisher"
)

func Pub(ctx context.Context, options *publisher.Options) (p publisher.Publisher, c func()) {
	p, err := publisher.New(options)
	if err != nil {
		log.Fatal(err)
	}

	startctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	if err := p.Start(startctx); err != nil {
		log.Fatal(err)
	}

	return p, func() {
		stopctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		if err := p.Stop(stopctx); err != nil {
			log.Println(err)
		}
	}
}
