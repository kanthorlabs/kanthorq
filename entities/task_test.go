package entities

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTaskId(t *testing.T) {
	require.True(t, strings.HasPrefix(TaskId(), "task"))
	require.True(t, strings.HasPrefix(TaskIdFromTime(time.Now()), "task"))
}
