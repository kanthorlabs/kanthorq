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
	if m == nil {
		return nil, errors.New("KANTHORQ.METADATA.VALUE.ERR: metadata is nil")

	}
	return json.Marshal(m)
}
