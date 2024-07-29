package kanthorq

import (
	"context"
)

var _ SubscriberPuller = (*DefaultSubscriberPuller)(nil)

// SubscriberPuller provider an abstraction of how to pull events from a stream
type SubscriberPuller interface {
	Pull(ctx context.Context) (*SubscriberPullerOut, error)
}

type SubscriberPullerOut struct {
	Tasks  map[string]*Task
	Events []*Event
}
