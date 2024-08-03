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

func Parse(id string) (string, ulid.ULID) {
	parts := strings.Split(id, Separator)
	if len(parts) != 2 {
		panic(fmt.Sprintf("[%s] does not start with namespace and end with ULID", id))
	}

	return parts[0], ulid.MustParse(parts[1])
}

func NewWithTime(ns string, t time.Time) string {
	return ns + Separator + strings.ToLower(ulid.MustNew(ulid.Timestamp(t), ulid.DefaultEntropy()).String())
}

// Next generates the next ID with the same namespace after a given duration
func Next(id string, duration time.Duration) string {
	ns, uid := Parse(id)

	future := ulid.Timestamp(ulid.Time(uid.Time()).Add(duration))
	if err := uid.SetTime(future); err != nil {
		panic(err)
	}

	return ns + Separator + strings.ToLower(uid.String())
}
