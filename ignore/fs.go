package ignore

import (
	"github.com/spf13/afero"
	"github.com/unmango/aferox/filter"
)

type Ignore interface {
	MatchesPath(string) bool
}

func NewFs(base afero.Fs, ignore Ignore) afero.Fs {
	return filter.FromPredicate(base, func(s string) bool {
		return !ignore.MatchesPath(s)
	})
}
