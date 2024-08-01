package kanthorq

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetadata_Scan(t *testing.T) {
	var m Metadata
	require.ErrorContains(t, m.Scan([]byte("ok")), "invalid character")
	require.ErrorContains(t, m.Scan(nil), "value is not []byte")
	require.NoError(t, m.Scan([]byte(`{"ok":true}`)))
	require.True(t, m["ok"].(bool))
}

func TestMetadata_Value(t *testing.T) {
	var nilm Metadata
	nilv, err := nilm.Value()
	require.NoError(t, err)
	require.NotNil(t, nilv)

	m := &Metadata{"ok": true}
	v, err := m.Value()
	require.NoError(t, err)
	require.NotNil(t, v)
}

func TestMetadata_Clone(t *testing.T) {
	m := Metadata{"ok": true}
	cm := m.Clone()
	cm["ok"] = false

	require.True(t, m["ok"].(bool))
	require.False(t, cm["ok"].(bool))
}

func TestMetadata_Merge(t *testing.T) {
	src := Metadata{"say": "hello"}
	m := Metadata{"ok": true}
	m.Merge(src)

	require.True(t, m["ok"].(bool))
	require.Equal(t, m["say"].(string), "hello")
}
