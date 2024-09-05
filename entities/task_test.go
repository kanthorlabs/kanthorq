package entities

import (
	"math"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const StateTesting TaskState = math.MaxInt16

func TestTaskId(t *testing.T) {
	require.True(t, strings.HasPrefix(TaskId(), "task"))
	require.True(t, strings.HasPrefix(TaskIdFromTime(time.Now()), "task"))
}

func TestTask_String(t *testing.T) {
	cases := map[TaskState]string{
		StateDiscarded: "discarded",
		StateCancelled: "cancelled",
		StatePending:   "pending",
		StateAvailable: "available",
		StateRunning:   "running",
		StateCompleted: "completed",
		StateRetryable: "retryable",
		StateTesting:   "",
	}

	for state, expected := range cases {
		require.Equal(t, expected, state.String())
	}
}

func TestAttemptedError_Scan(t *testing.T) {
	var m AttemptedError

	require.ErrorContains(t, m.Scan([]byte("ok")), "invalid character")
	require.ErrorContains(t, m.Scan(nil), "only string or []byte supported")

	require.NoError(t, m.Scan([]byte(`{"error":"some"}`)))
	require.Equal(t, m.Error, "some")

	require.NoError(t, m.Scan(`{"error":"some"}`))
	require.Equal(t, m.Error, "some")
}

func TestAttemptedError_Value(t *testing.T) {
	var nilm AttemptedError
	nilv, err := nilm.Value()
	require.NoError(t, err)
	require.NotNil(t, nilv)

	m := &AttemptedError{Error: "some"}
	v, err := m.Value()
	require.NoError(t, err)
	require.NotNil(t, v)
}
