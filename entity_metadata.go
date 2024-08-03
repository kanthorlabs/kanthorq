package kanthorq

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Metadata map[string]interface{}

// Scan implements the sql.Scanner interface to scan a value from the database into the Metadata struct
func (m *Metadata) Scan(value interface{}) error {
	if data, ok := value.([]byte); ok {
		return json.Unmarshal(data, m)
	}

	if data, ok := value.(string); ok {
		return json.Unmarshal([]byte(data), m)
	}

	return fmt.Errorf("KANTHORQ.METADATA.SCAN.ERR: only string or []byte supported, got %T", value)
}

// Value implements the driver.Valuer interface to convert the Metadata struct to a value that can be stored in the database
func (m Metadata) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m Metadata) Clone() Metadata {
	c := Metadata{}
	for k, v := range m {
		c[k] = v
	}
	return c
}

func (m Metadata) Merge(src Metadata) {
	for k, v := range src {
		m[k] = v
	}
}
