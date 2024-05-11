package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func AbsPathify(in string) string {
	if in == "$HOME" || strings.HasPrefix(in, "$HOME"+string(os.PathSeparator)) {
		in = os.Getenv("HOME") + in[5:]
	}

	in = os.ExpandEnv(in)

	if filepath.IsAbs(in) {
		return filepath.Clean(in)
	}

	p, err := filepath.Abs(in)
	if err == nil {
		return filepath.Clean(p)
	}

	return ""
}
