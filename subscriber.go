package kanthorq

import (
	"context"
	"errors"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/pkg/validator"
)

type Subscriber interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Receive(ctx context.Context, handler SubscriberHandler) error
}

type SubscriberHandler func(ctx context.Context, event *Event) error

func NewSubscriber(uri string, options *SubscriberOptions) (Subscriber, error) {
	if err := validator.Validate.Struct(options); err != nil {
		return nil, err
	}
	return &subscriber{uri: uri, options: options}, nil
}

type subscriber struct {
	uri     string
	options *SubscriberOptions
	mu      sync.Mutex

	conn     *pgx.Conn
	stream   *StreamRegistry
	consumer *ConsumerRegistry
}

func (pub *subscriber) Start(ctx context.Context) error {
	if err := pub.connect(ctx); err != nil {
		return err
	}

	pub.mu.Lock()
	defer pub.mu.Unlock()

	tx, err := pub.conn.Begin(ctx)
	if err != nil {
		return err
	}
	req := &ConsumerRegisterReq{
		StreamName:         pub.options.StreamName,
		ConsumerName:       pub.options.ConsumerName,
		ConsumerTopic:      pub.options.ConsumerTopic,
		ConsumerAttemptMax: pub.options.ConsumerAttemptMax,
	}
	res, err := req.Do(ctx, tx)
	if err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	pub.stream = res.StreamRegistry
	pub.consumer = res.ConsumerRegistry

	return nil
}

func (pub *subscriber) Stop(ctx context.Context) (err error) {
	pub.mu.Lock()
	defer pub.mu.Unlock()

	if pub.conn != nil {
		err = errors.Join(err, pub.conn.Close(ctx))
	}

	pub.conn = nil
	pub.stream = nil
	pub.consumer = nil
	return
}

func (pub *subscriber) connect(ctx context.Context) error {
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

func (pub *subscriber) Receive(ctx context.Context, handler SubscriberHandler) error {
	return errors.New("unimplemented")
}
