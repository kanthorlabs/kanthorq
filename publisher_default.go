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

	conn, err := pub.cm.Acquire(ctx)
	if err != nil {
		return
	}
	defer pub.cm.Release(ctx, conn)

	req := &StreamRegisterReq{StreamName: pub.options.StreamName}
	res, err := Do(ctx, req, conn)
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

func (pub *publisher) Send(ctx context.Context, events ...*Event) error {
	if len(events) == 0 {
		return errors.New("no events provided")
	}

	for i, e := range events {
		if err := validator.Validate.Struct(e); err != nil {
			return fmt.Errorf("event %d: %w", i, err)
		}
	}

	conn, err := pub.cm.Acquire(ctx)
	if err != nil {
		return err
	}
	defer pub.cm.Release(ctx, conn)

	req := &StreamPutEventsReq{Stream: pub.stream, Events: events}
	res, err := Do(ctx, req, conn)
	if err != nil {
		return err
	}

	if res.InsertCount != int64(len(events)) {
		log.Println("inserted ", res.InsertCount, " events, expected ", len(events))
	}

	return nil
}
