package entities

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConsumerId(t *testing.T) {
	require.True(t, strings.HasPrefix(ConsumerId(), "consumer"))
	require.True(t, strings.HasPrefix(ConsumerIdFromTime(time.Now()), "consumer"))
}
