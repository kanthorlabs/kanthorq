package kanthorq

import (
	"context"
	"errors"
	"sync"

	"github.com/jackc/pgx/v5"
)

var _ Publisher = (*publisher)(nil)

type Publisher interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Send(ctx context.Context, events ...*Event) error
}

// NewPublisher creates a new publisher that uses the default stream
func NewPublisher(uri string, options *PublisherOptions) (Publisher, error) {
	if err := options.Validate(); err != nil {
		return nil, err
	}
	return &publisher{uri: uri, options: options}, nil
}

type publisher struct {
	uri     string
	options *PublisherOptions
	mu      sync.Mutex

	conn   *pgx.Conn
	stream *StreamRegistry
}

func (pub *publisher) Start(ctx context.Context) error {
	if err := pub.connect(ctx); err != nil {
		return err
	}

	pub.mu.Lock()
	defer pub.mu.Unlock()

	tx, err := pub.conn.Begin(ctx)
	if err != nil {
		return err
	}
	req := &StreamRegisterReq{StreamName: pub.options.StreamName}
	res, err := req.Do(ctx, tx)
	if err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	pub.stream = res.StreamRegistry

	return nil
}

func (pub *publisher) Stop(ctx context.Context) (err error) {
	pub.mu.Lock()
	defer pub.mu.Unlock()

	if pub.conn != nil {
		err = errors.Join(err, pub.conn.Close(ctx))
	}

	pub.conn = nil
	pub.stream = nil
	return
}

func (pub *publisher) connect(ctx context.Context) error {
	pub.mu.Lock()
	defer pub.mu.Unlock()

	// connection is already ready, don't need to re-connect
	if pub.conn != nil && !pub.conn.IsClosed() {
		return nil
	}

	conn, err := pgx.Connect(ctx, pub.uri)
	if err != nil {
		return err
	}
	pub.conn = conn

	return nil
}

func (pub *publisher) Send(ctx context.Context, events ...*Event) error {
	// @TODO: validate events

	// @TODO: throttle sending events
	// @TODO: circuit breaker
	// @TODO: retry logic

	tx, err := pub.conn.Begin(ctx)
	if err != nil {
		return err
	}
	resp, err := StreamPut(pub.stream, events).Do(ctx, tx)
	if err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}

	if resp.InsertCount != int64(len(events)) {
		// TODO: log error about inserted count and expected count does not match
	}

	return nil
}