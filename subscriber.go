package kanthorq

import (
	"context"
	"errors"
	"sync"

	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
	"github.com/kanthorlabs/kanthorq/pkg/validator"
)

type Subscriber interface {
	Start(ctx context.Context) (err error)
	Stop(ctx context.Context) (err error)
	Receive(ctx context.Context, handler SubscriberHandler) (err error)
}

type SubscriberHandler func(ctx context.Context, event *Event) error

func NewSubscriber(uri string, options *SubscriberOptions) (Subscriber, error) {
	if err := validator.Validate.Struct(options); err != nil {
		return nil, err
	}
	cm, err := pgcm.New(uri)
	if err != nil {
		return nil, err
	}
	return &subscriber{cm: cm, options: options}, nil
}

type subscriber struct {
	cm      pgcm.ConnectionManager
	options *SubscriberOptions
	mu      sync.Mutex

	stream   *StreamRegistry
	consumer *ConsumerRegistry
}

func (pub *subscriber) Start(ctx context.Context) (err error) {
	pub.mu.Lock()
	defer pub.mu.Unlock()

	if err = pub.cm.Start(ctx); err != nil {
		return
	}

	conn, err := pub.cm.Connection(ctx)
	if err != nil {
		return
	}
	defer func() { err = conn.Close(ctx) }()

	req := &ConsumerRegisterReq{
		StreamName:         pub.options.StreamName,
		ConsumerName:       pub.options.ConsumerName,
		ConsumerTopic:      pub.options.ConsumerTopic,
		ConsumerAttemptMax: pub.options.ConsumerAttemptMax,
	}
	res, err := Do(ctx, req, conn.Raw())
	if err != nil {
		return
	}

	pub.stream = res.StreamRegistry
	pub.consumer = res.ConsumerRegistry
	return nil
}

func (pub *subscriber) Stop(ctx context.Context) (err error) {
	pub.mu.Lock()
	defer pub.mu.Unlock()

	if cmerr := pub.cm.Start(ctx); cmerr != nil {
		err = errors.Join(err, cmerr)
	}

	pub.stream = nil
	pub.consumer = nil
	return
}

func (pub *subscriber) Receive(ctx context.Context, handler SubscriberHandler) error {
	return errors.New("unimplemented")
}
