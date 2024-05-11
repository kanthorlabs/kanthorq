package safe

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

// Metadata is a thread-safe key-value store.
// There are some gotchas when using this struct:
// json marshalling is supported but all number will be parsed as float64.
// yaml marshalling is supported but all number will be parsed as int.
type Metadata struct {
	kv map[string]any
	mu sync.Mutex
}

func (meta *Metadata) Set(k string, v any) {
	meta.mu.Lock()
	defer meta.mu.Unlock()

	if meta.kv == nil {
		meta.kv = make(map[string]any)
	}

	meta.kv[k] = v
}

func (meta *Metadata) Get(k string) (any, bool) {
	meta.mu.Lock()
	defer meta.mu.Unlock()

	if meta.kv == nil {
		return nil, false
	}

	v, has := meta.kv[k]
	return v, has
}

func (meta *Metadata) Merge(src *Metadata) {
	meta.mu.Lock()
	defer meta.mu.Unlock()

	if meta.kv == nil {
		meta.kv = make(map[string]any)
	}

	if src == nil || len(src.kv) == 0 {
		return
	}

	for k := range src.kv {
		meta.kv[k] = src.kv[k]
	}
}

func (meta *Metadata) String() string {
	if meta == nil || meta.kv == nil {
		return ""
	}

	data, _ := json.Marshal(meta.kv)
	return string(data)
}

// Value implements the driver Valuer interface.
func (meta *Metadata) Value() (driver.Value, error) {
	// meta == nil when we convert it to database value
	if meta == nil || meta.kv == nil {
		return "", nil
	}
	data, err := json.Marshal(meta.kv)
	return string(data), err
}

// Scan implements the Scanner interface.
func (meta *Metadata) Scan(value any) error {
	v := value.(string)
	if v == "" {
		return nil
	}
	return json.Unmarshal([]byte(v), &meta.kv)
}

func (meta *Metadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(meta.kv)
}

func (meta *Metadata) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &meta.kv)
}

func (meta *Metadata) MarshalYAML() (interface{}, error) {
	var value yaml.Node
	return value, value.Encode(meta.kv)
}

func (meta *Metadata) UnmarshalYAML(value *yaml.Node) error {
	return value.Decode(&meta.kv)
}

func (meta *Metadata) FromHttpHeader(headers http.Header) {
	for k, v := range headers {
		meta.Set(strings.ToLower(k), v[0])
	}
}

func (meta *Metadata) ToHttpHeader() http.Header {
	headers := http.Header{}
	for k, v := range meta.kv {
		switch value := v.(type) {
		case string:
			headers.Set(k, value)
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			headers.Set(k, fmt.Sprintf("%d", value))
		case float32, float64:
			headers.Set(k, fmt.Sprintf("%f", value))
		case bool:
			headers.Set(k, fmt.Sprintf("%t", value))
			// ignore other cases
		}
	}
	return headers
}

func MetadataMapstructureHook() mapstructure.DecodeHookFuncType {
	return func(from, to reflect.Type, data interface{}) (interface{}, error) {
		if from.Kind() == reflect.Map && to.String() == "*safe.Metadata" {
			var metdata Metadata
			for k, v := range data.(map[string]interface{}) {
				metdata.Set(k, v)
			}
			return &metdata, nil
		}

		return data, nil
	}
}
