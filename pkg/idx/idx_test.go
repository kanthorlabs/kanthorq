package idx

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIdx_New(t *testing.T) {
	ns := "kanthorq"
	id := New(ns)
	require.True(t, strings.HasPrefix(id, ns))
}
