package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatchSubject(t *testing.T) {
	cases := []struct {
		pattern string
		subject string
		matched bool
	}{
		// base cases
		{"", "time.us.east", false},
		{"time.*.east", "", false},
		{"*", "", false},
		// case sensitive
		{"time.us.east", "time.us.EAST", false},

		{"time.*.east", "time.us.east", true},
		{"time.*.east", "time.eu.east", true},
		{"time.*.east", "time.NewYork.east", true},
		{"time.us.*", "time.us.east", true},
		{"time.us.*", "time.us.east.atlanta", false},

		{"time.>.east", "time.us.east", false},
		{"time.us.>", "time.us.east", true},
		{"time.us.>", "time.us.east.atlanta", true},

		// mixed case
		{"*.*.east.>", "time.asia.east", false},
		{"*.*.east.>", "time.us.east.atlanta", true},
		{"*.*.east.>", "time.eu.east.london", true},

		{"time.*.east.*", "time.us.west.any", false},
		{"time.us.*.east", "time.us.any.east", true},
		{"time.*.east.*", "time.us.east.any", true},

		{"time.us.east.>", "time.us.east", false},
		{"time.us.east.>", "time.us.east.newyork", true},
		{"time.us.east.>", "time.us.west", false},

		// mixed case
		{"time.*.east.>", "time.us.east.atlanta", true},
		{"time.*.east.>", "time.us.west", false},
		{"time.*.east.>", "time.us.east", false},
	}

	for _, c := range cases {
		matched := MatchSubject(c.pattern, c.subject)
		msg := fmt.Sprintf("%s -> %s | expected:%v # actual:%v", c.pattern, c.subject, c.matched, matched)
		require.Equal(t, matched, c.matched, msg)
	}
}
