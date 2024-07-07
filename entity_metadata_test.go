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
	_, err := nilm.Value()
	require.ErrorContains(t, err, "metadata is ni")

	m := Metadata{"ok": true}
	_, err = m.Value()
	require.NoError(t, err)
}
