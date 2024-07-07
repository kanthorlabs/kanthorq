package idx

import (
	"fmt"

	"github.com/oklog/ulid/v2"
)

// New generates an ID that is Lexicographically Sortable with the given namespace
func New(ns string) string {
	return fmt.Sprintf("%s_%s", ns, ulid.Make().String())
}
