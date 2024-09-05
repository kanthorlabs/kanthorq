package entities

import (
	"strings"
	"testing"
	"time"

	"github.com/kanthorlabs/kanthorq/pkg/xvalidator"
	"github.com/stretchr/testify/require"
)

func TestEventId(t *testing.T) {
	require.True(t, strings.HasPrefix(EventId(), "event"))
	require.True(t, strings.HasPrefix(EventIdFromTime(time.Now()), "event"))
}

func TestNewEvent(t *testing.T) {
	require.NoError(t, xvalidator.Validate.Struct(NewEvent("ok", []byte("ok"))))
}
