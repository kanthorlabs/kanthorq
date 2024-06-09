package subscriber

import "time"

type Options struct {
	// How many events to consume at batch
	Size int
	// a period of time during which KanthorQ prevents
	// all consumers from receiving and processing the message.
	VisibilityTimeout time.Duration
	// a period of time during which KanthorQ waits for
	// if there is no message in current batch
	WaitingTime time.Duration
}

type Option func(options *Options)

var DefaultSize = 100

func Size(size int) Option {
	return func(options *Options) {
		options.Size = size
	}
}

var DefaultVisibilityTimeout = time.Minute * 15

func VisibilityTimeout(vt time.Duration) Option {
	return func(options *Options) {
		options.VisibilityTimeout = vt
	}
}

var DefaultWaitingTime = time.Second * 15

func WaitingTime(wt time.Duration) Option {
	return func(options *Options) {
		options.WaitingTime = wt
	}
}
