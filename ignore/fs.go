package ignore

import (
	"github.com/spf13/afero"
	"github.com/unmango/aferox/filter"
	"github.com/unmango/aferox/op"
)

type Ignore interface {
	MatchesPath(string) bool
}

func NewFs(base afero.Fs, ignore Ignore) afero.Fs {
	return filter.FromPredicate(base, func(op op.Operation) bool {
		return !ignore.MatchesPath(op.Path())
	})
}
