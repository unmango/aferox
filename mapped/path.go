package mapped

import (
	"path"
	"strings"
)

func CutPrefix(s, prefix string) (after string, found bool) {
	if after, found = strings.CutPrefix(path.Clean(s), prefix); found {
		return path.Clean(after), found
	}

	return
}
