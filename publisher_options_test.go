package kanthorq

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPublisherOptions_Validate(t *testing.T) {
	require.NotNil(t, (&PublisherOptions{}).Validate())

	options := &PublisherOptions{StreamName: DefaultStreamName}
	require.NoError(t, options.Validate())
}
