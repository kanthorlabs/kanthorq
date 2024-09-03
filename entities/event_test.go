package entities

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEventId(t *testing.T) {
	require.True(t, strings.HasPrefix(EventId(), "event"))
	require.True(t, strings.HasPrefix(EventIdFromTime(time.Now()), "event"))
}
