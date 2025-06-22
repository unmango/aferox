package mapped

import (
	"strings"
)

func CutPrefix(s, prefix string) (after string, found bool) {
	prefix = strings.TrimLeft(prefix, "/")
	s, abs := strings.CutPrefix(s, "/")

	if after, found = strings.CutPrefix(s, prefix); !found {
		return
	}

	if abs {
		return after, true
	} else {
		return strings.TrimLeft(after, "/"), true
	}
}
