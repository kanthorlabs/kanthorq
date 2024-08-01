package kanthorq

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Metadata map[string]interface{}

// Scan implements the sql.Scanner interface to scan a value from the database into the Metadata struct
func (m *Metadata) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("KANTHORQ.METADATA.SCAN.ERR: value is not []byte")
	}

	return json.Unmarshal(bytes, m)
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
