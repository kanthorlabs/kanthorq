package subscriber

import (
	"context"

	"github.com/kanthorlabs/kanthorq/pkg/xvalidator"
	"github.com/kanthorlabs/kanthorq/puller"
	"go.uber.org/zap"
)

func New(options *Options, logger *zap.Logger) (Subscriber, error) {
	if err := xvalidator.Validate.Struct(options); err != nil {
		return nil, err
	}

	logger = logger.With(
		zap.String("subscriber", "primary"),
		zap.String("puller", "primary"),
		zap.String("stream_name", options.StreamName),
		zap.String("consumer_name", options.ConsumerName),
	)
	return &primary{options: options, logger: logger, pullerF: puller.New}, nil
}

func NewRetry(options *Options, logger *zap.Logger) (Subscriber, error) {
	if err := xvalidator.Validate.Struct(options); err != nil {
		return nil, err
	}

	logger = logger.With(
		zap.String("subscriber", "primary"),
		zap.String("puller", "retry"),
		zap.String("stream_name", options.StreamName),
		zap.String("consumer_name", options.ConsumerName),
	)
	return &primary{options: options, logger: logger, pullerF: puller.NewRetry}, nil
}

func NewVisibility(options *Options, logger *zap.Logger) (Subscriber, error) {
	if err := xvalidator.Validate.Struct(options); err != nil {
		return nil, err
	}

	logger = logger.With(
		zap.String("subscriber", "primary"),
		zap.String("puller", "visibility"),
		zap.String("stream_name", options.StreamName),
		zap.String("consumer_name", options.ConsumerName),
	)
	return &primary{options: options, logger: logger, pullerF: puller.NewVisibility}, nil
}

type Subscriber interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Receive(ctx context.Context, handler Handler) error
}

type Handler func(ctx context.Context, msg *Message) error
