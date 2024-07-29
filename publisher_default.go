package kanthorq

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
	"github.com/kanthorlabs/kanthorq/pkg/validator"
)

var _ Publisher = (*publisher)(nil)

type publisher struct {
	options *PublisherOptions
	mu      sync.Mutex

	cm     pgcm.ConnectionManager
	stream *StreamRegistry
}

func (pub *publisher) Start(ctx context.Context) (err error) {
	pub.mu.Lock()
	defer pub.mu.Unlock()

	if err = pub.cm.Start(ctx); err != nil {
		return
	}

	conn, err := pub.cm.Connection(ctx)
	if err != nil {
		return
	}
	defer conn.Close(ctx)

	req := &StreamRegisterReq{StreamName: pub.options.StreamName}
	res, err := Do(ctx, req, conn.Raw())
	if err != nil {
		return
	}

	pub.stream = res.StreamRegistry
	return nil
}

func (pub *publisher) Stop(ctx context.Context) (err error) {
	pub.mu.Lock()
	defer pub.mu.Unlock()

	if cmerr := pub.cm.Start(ctx); cmerr != nil {
		err = errors.Join(err, cmerr)
	}

	pub.stream = nil
	return
}

func (pub *publisher) Send(ctx context.Context, events ...*Event) (err error) {
	if len(events) == 0 {
		err = errors.New("no events provided")
		return
	}

	for i, e := range events {
		if err = validator.Validate.Struct(e); err != nil {
			err = fmt.Errorf("event %d: %w", i, err)
			return
		}
	}

	conn, err := pub.cm.Connection(ctx)
	if err != nil {
		return
	}
	defer conn.Close(ctx)

	req := &StreamPutEventsReq{Stream: pub.stream, Events: events}
	res, err := Do(ctx, req, conn.Raw())
	if err != nil {
		return
	}

	if res.InsertCount != int64(len(events)) {
		log.Println("inserted ", res.InsertCount, " events, expected ", len(events))
	}

	return nil
}
