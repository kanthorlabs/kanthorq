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
	"go.uber.org/zap"
)

var _ Publisher = (*primary)(nil)

// NewPublisher creates a new publisher that uses the default stream
func New(options *Options, logger *zap.Logger) (Publisher, error) {
	if err := xvalidator.Validate.Struct(options); err != nil {
		return nil, err
	}

	logger = logger.With(
		zap.String("publisher", "primary"),
		zap.String("stream_name", options.StreamName),
	)
	return &primary{options: options, logger: logger}, nil
}

type primary struct {
	mu sync.Mutex

	options *Options
	logger  *zap.Logger
	cm      pgcm.ConnectionManager

	stream *entities.StreamRegistry
}

func (pub *primary) Start(ctx context.Context) error {
	pub.mu.Lock()
	defer pub.mu.Unlock()

	cm, err := pgcm.New(pub.options.Connection)
	if err != nil {
		return err
	}
	pub.cm = cm
	if err = pub.cm.Start(ctx); err != nil {
		return err
	}

	conn, err := pub.cm.Acquire(ctx)
	if err != nil {
		return err
	}
	defer pub.cm.Release(ctx, conn)

	req := &core.StreamRegisterReq{StreamName: pub.options.StreamName}
	res, err := core.Do(ctx, conn, req)
	if err != nil {
		return err
	}

	pub.stream = res.StreamRegistry
	pub.logger.Info("started")
	return nil
}

func (pub *primary) Stop(ctx context.Context) error {
	pub.mu.Lock()
	defer pub.mu.Unlock()

	var err error
	if cmerr := pub.cm.Stop(ctx); cmerr != nil {
		err = errors.Join(err, cmerr)
	}

	pub.stream = nil
	pub.logger.Info("stopped")
	return err
}

func (pub *primary) Send(ctx context.Context, events []*entities.Event) error {
	if len(events) == 0 {
		return errors.New("PUBLISHER.SEND.NO_EVENTS")
	}

	for i, e := range events {
		if err := xvalidator.Validate.Struct(e); err != nil {
			return fmt.Errorf("PUBLISHER.SEND.EVENT.%d: %w", i, err)
		}
	}

	req := &core.StreamPutEventsReq{Stream: pub.stream, Events: events}
	res, err := core.DoWithCM(ctx, pub.cm, req)
	if err != nil {
		return err
	}

	pub.logger.Info("sent events", zap.Int("event_count", len(events)), zap.Int64("insert_count", res.InsertCount))
	return nil
}

func (pub *primary) SendTx(ctx context.Context, events []*entities.Event, tx pgx.Tx) error {
	req := &core.StreamPutEventsReq{Stream: pub.stream, Events: events}
	_, err := req.Do(ctx, tx)
	return err
}
