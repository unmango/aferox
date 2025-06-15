package mapped

import (
	"path"
	"strings"
)

func Clean(p string) string {
	return strings.TrimLeft(path.Clean(p), "/")
}

func CutPrefix(s, prefix string) (after string, found bool) {
	if after, found = strings.CutPrefix(Clean(s), prefix); found {
		return Clean(after), found
	}

	return
}
