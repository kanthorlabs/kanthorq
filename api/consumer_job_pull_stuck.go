package api

import (
	_ "embed"
	"time"

	"github.com/kanthorlabs/kanthorq/entities"
)

// NewConsumerJobPullStuck will find retryable jobs, make them become running, and return the events themself
func NewConsumerJobPullStuck(consumer *entities.Consumer, size int, vt time.Duration) *ConsumerJobPullReq {
	return &ConsumerJobPullReq{
		Consumer:          consumer,
		Size:              size,
		VisibilityTimeout: vt,
		FromState:         entities.StateRunning,
		ToState:           entities.StateRunning,
		Source:            "ConsumerJobPullStuck",
	}
}
