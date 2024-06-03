package subscriber

import "time"

type Options struct {
	Size              int
	VisibilityTimeout time.Duration
}

type Option func(options *Options)

var DefaultSize = 100

func Size(size int) Option {
	return func(options *Options) {
		options.Size = size
	}
}

var DefaultVisibilityTimeout = time.Hour

func VisibilityTimeout(vt time.Duration) Option {
	return func(options *Options) {
		options.VisibilityTimeout = vt
	}
}
