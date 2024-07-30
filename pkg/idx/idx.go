package idx

import (
	"fmt"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
)

var Separator = "_"

// New generates an ID that is Lexicographically Sortable with the given namespace
func New(ns string) string {
	return ns + Separator + strings.ToLower(ulid.Make().String())
}

// Next generates the next ID with the same namespace after a given duration
func Next(id string, duration time.Duration) string {
	parts := strings.Split(id, Separator)
	if len(parts) != 2 {
		panic(fmt.Sprintf("[%s] does not start with namespace and end with ULID", id))
	}

	uid := ulid.MustParse(parts[1])
	upper := ulid.Time(uid.Time() + uint64(duration.Milliseconds()))
	uid.SetTime(ulid.Timestamp(upper))

	return parts[0] + Separator + strings.ToLower(uid.String())
}
