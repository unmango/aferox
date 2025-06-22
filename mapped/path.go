package mapped

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
)

func CutPrefix(s, prefix string) (after string, found bool) {
	if after, found = strings.CutPrefix(s, prefix); !found {
		return
	}

	if path.IsAbs(s) {
		return fmt.Sprint(filepath.Separator, after), true
	} else {
		return after, true
	}
}
