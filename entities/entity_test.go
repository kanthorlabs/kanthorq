package entities

import (
	"strings"
	"testing"

	"github.com/kanthorlabs/kanthorq/pkg/xfaker"
	"github.com/stretchr/testify/require"
)

func TestCollection(t *testing.T) {
	require.True(t, strings.HasPrefix(Collection(xfaker.StreamName()), "kanthorq"))
}

func TestProperties(t *testing.T) {
	var value = struct {
		Name string `json:"name"`
		Id   string `json:"id"`
		Age  int    `json:"age"`
	}{}

	require.Equal(t, []string{"name", "id", "age"}, Properties(value))
	require.Equal(t, []string{"name", "id", "age"}, Properties(&value))
}
