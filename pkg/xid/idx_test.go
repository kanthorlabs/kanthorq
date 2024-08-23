package xid

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIdx_New(t *testing.T) {
	ns := "kanthorq"
	id := New(ns)
	require.True(t, strings.HasPrefix(id, ns))
}

func TestIdx_Parse(t *testing.T) {
	ns := "kanthorq"
	id := New(ns)

	parsedNs, parsedId := Parse(id)
	require.Equal(t, ns, parsedNs)
	require.NotNil(t, id, parsedId)
}

func TestIdx_NewWithTime(t *testing.T) {
	ns := "kanthorq"

	now := time.Now()
	future := now.Add(time.Hour)

	id := NewWithTime(ns, future)
	require.True(t, strings.HasPrefix(id, ns))

	_, parsedId := Parse(id)

	require.Equal(t, future.UnixMilli(), int64(parsedId.Time()))
}

func TestIdx_Next(t *testing.T) {
	ns := "kanthorq"
	duration := time.Hour

	now := time.Now()
	start := NewWithTime(ns, now)

	end := Next(start, duration)
	require.True(t, strings.HasPrefix(end, ns))

	require.True(t, start < end)

	_, parsedId := Parse(end)
	require.Equal(t, now.Add(duration).UnixMilli(), int64(parsedId.Time()))
}
