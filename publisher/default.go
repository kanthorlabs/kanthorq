package publisher

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/kanthorlabs/kanthorq/core"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
	"github.com/kanthorlabs/kanthorq/pkg/validator"
)

var _ Publisher = (*primary)(nil)

type primary struct {
	options *Options
	mu      sync.Mutex

	cm     pgcm.ConnectionManager
	stream *entities.StreamRegistry
}

func (pub *primary) Start(ctx context.Context) (err error) {
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

	req := &core.StreamRegisterReq{StreamName: pub.options.StreamName}
	res, err := core.Do(ctx, req, conn)
	if err != nil {
		return
	}

	pub.stream = res.StreamRegistry
	return nil
}

func (pub *primary) Stop(ctx context.Context) (err error) {
	pub.mu.Lock()
	defer pub.mu.Unlock()

	if cmerr := pub.cm.Start(ctx); cmerr != nil {
		err = errors.Join(err, cmerr)
	}

	pub.stream = nil
	return
}

func (pub *primary) Send(ctx context.Context, events ...*entities.Event) error {
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

	req := &core.StreamPutEventsReq{Stream: pub.stream, Events: events}
	res, err := core.Do(ctx, req, conn)
	if err != nil {
		return err
	}

	if res.InsertCount != int64(len(events)) {
		log.Println("inserted ", res.InsertCount, " events, expected ", len(events))
	}

	return nil
}
