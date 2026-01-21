package filter_test

import (
	"github.com/spf13/afero"
	"github.com/unmango/aferox/filter"
)

// Helper functions for backward compatibility during migration
func newFsWithPathPredicate(base afero.Fs, pred func(string) bool) afero.Fs {
	return filter.NewFs(base, filter.PathPredicate(pred))
}

func fromPredicateWithErrorAndPath(base afero.Fs, pred func(string) bool, err error) afero.Fs {
	return filter.FromPredicateWithError(base, filter.PathPredicate(pred), err)
}
