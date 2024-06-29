package q

import (
	_ "embed"
	"time"

	"github.com/kanthorlabs/kanthorq/entities"
)

// NewConsumerJobPullAvailable is the main function of this system
// it will find available jobs, make it become running, and return the events themself
func NewConsumerJobPullAvailable(consumer *entities.Consumer, size int, vt time.Duration) *ConsumerJobPullReq {
	return &ConsumerJobPullReq{
		Consumer:          consumer,
		Size:              size,
		VisibilityTimeout: vt,
		FromState:         entities.StateAvailable,
		ToState:           entities.StateRunning,
		Source:            "ConsumerJobPullAvailable",
	}
}
