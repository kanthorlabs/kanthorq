package subscriber

import "time"

func NewSubscribeOption(options ...SubscribeOption) *SubscribeOptions {
	opts := &SubscribeOptions{
		Size:              DefaultSize,
		Timeout:           DefaultTimeout,
		VisibilityTimeout: DefaultVisibilityTimeout,
		WaitingTime:       DefaultWaitingTime,
	}
	for _, configure := range options {
		configure(opts)
	}
	return opts
}

type SubscribeOptions struct {
	// How many events to consume at batch
	Size int
	// a period of time during which KanthorQ waits for pulling events
	Timeout time.Duration
	// a period of time during which KanthorQ prevents
	// all consumers from receiving a message after pulling.
	VisibilityTimeout time.Duration
	// a period of time during which KanthorQ waits for
	// if there is no message in current batch
	WaitingTime time.Duration
}

type SubscribeOption func(options *SubscribeOptions)

var DefaultSize = 100

func Size(size int) SubscribeOption {
	return func(options *SubscribeOptions) {
		options.Size = size
	}
}

var DefaultTimeout = time.Minute

func Timeout(vt time.Duration) SubscribeOption {
	return func(options *SubscribeOptions) {
		options.Timeout = vt
	}
}

var DefaultVisibilityTimeout = time.Minute * 15

func VisibilityTimeout(vt time.Duration) SubscribeOption {
	return func(options *SubscribeOptions) {
		options.VisibilityTimeout = vt
	}
}

var DefaultWaitingTime = time.Second * 5

func WaitingTime(wt time.Duration) SubscribeOption {
	return func(options *SubscribeOptions) {
		options.WaitingTime = wt
	}
}
