package kanthorq

import (
	"context"
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
func NewPublisher(uri string) (Publisher, error) {
	return &publisher{uri: uri, stream: &StreamRegistry{Name: DefaultStreamName}}, nil
}

type publisher struct {
	uri string
	mu  sync.Mutex

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
	resp, err := StreamRegister(pub.stream.Name).Do(ctx, tx)
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	pub.stream = resp.StreamRegistry

	return nil
}

func (pub *publisher) Stop(ctx context.Context) error {
	pub.mu.Lock()
	defer pub.mu.Unlock()

	pub.stream = nil

	if pub.conn == nil {
		return nil
	}
	if err := pub.conn.Close(ctx); err != nil {
		return err
	}

	pub.conn = nil
	return nil
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
