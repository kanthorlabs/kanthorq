package publisher

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/core"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
	"github.com/kanthorlabs/kanthorq/pkg/xvalidator"
)

var _ Publisher = (*primary)(nil)

// NewPublisher creates a new publisher that uses the default stream
func New(options *Options) (Publisher, error) {
	if err := xvalidator.Validate.Struct(options); err != nil {
		return nil, err
	}
	cm, err := pgcm.New(options.Connection)
	if err != nil {
		return nil, err
	}
	return &primary{options: options, cm: cm}, nil
}

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

	if cmerr := pub.cm.Stop(ctx); cmerr != nil {
		err = errors.Join(err, cmerr)
	}

	pub.stream = nil
	return
}

func (pub *primary) Send(ctx context.Context, events []*entities.Event) error {
	if len(events) == 0 {
		return errors.New("no events provided")
	}

	for i, e := range events {
		if err := xvalidator.Validate.Struct(e); err != nil {
			return fmt.Errorf("event %d: %w", i, err)
		}
	}

	req := &core.StreamPutEventsReq{Stream: pub.stream, Events: events}
	_, err := core.DoWithCM(ctx, req, pub.cm)
	return err
}

func (pub *primary) SendTx(ctx context.Context, events []*entities.Event, tx pgx.Tx) error {
	req := &core.StreamPutEventsReq{Stream: pub.stream, Events: events}
	_, err := req.Do(ctx, tx)
	return err
}
