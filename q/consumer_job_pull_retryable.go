package q

import (
	_ "embed"
	"time"

	"github.com/kanthorlabs/kanthorq/entities"
)

// NewConsumerJobPullRetryable will find retryable jobs, make them become running, and return the events themself
func NewConsumerJobPullRetryable(consumer *entities.Consumer, size int, vt time.Duration) *ConsumerJobPullReq {
	return &ConsumerJobPullReq{
		Consumer:          consumer,
		Size:              size,
		VisibilityTimeout: vt,
		FromState:         entities.StateRetryable,
		ToState:           entities.StateRunning,
		Source:            "ConsumerJobPullRetryable",
	}
}
