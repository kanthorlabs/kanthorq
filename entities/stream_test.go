package entities

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStreamId(t *testing.T) {
	require.True(t, strings.HasPrefix(StreamId(), "stream"))
	require.True(t, strings.HasPrefix(StreamIdFromTime(time.Now()), "stream"))
}
