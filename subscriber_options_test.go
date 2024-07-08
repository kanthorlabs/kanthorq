package kanthorq

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubscriberOptions(t *testing.T) {
	require.NotNil(t, (&SubscriberOptions{}).Validate())

	options := &SubscriberOptions{
		StreamName:         DefaultStreamName,
		ConsumerName:       "internal",
		ConsumerTopic:      "system.ping",
		ConsumerAttemptMax: DefaultConsumerAttemptMax,
	}
	require.NoError(t, options.Validate())
}
