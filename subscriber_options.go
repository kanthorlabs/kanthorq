package kanthorq

import "time"

type SubscriberOptions struct {
	Size              int
	VisibilityTimeout time.Duration
}

type SubscribeOption func(options *SubscriberOptions)

var DefaultSubscriberSize = 100

func SubscriberSize(size int) SubscribeOption {
	return func(options *SubscriberOptions) {
		options.Size = size
	}
}

var DefaultSubscriberVisibilityTimeout = time.Hour

func SubscriberVisibilityTimeout(vt time.Duration) SubscribeOption {
	return func(options *SubscriberOptions) {
		options.VisibilityTimeout = vt
	}
}
