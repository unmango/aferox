package aferox

import (
	"io/fs"

	"github.com/spf13/afero"
	"github.com/unmango/go/iter"
	"github.com/unmango/go/option"
)

type ErrFilter func(error) error

type iterOptions struct {
	continueOnErr bool
	errFilter     ErrFilter
	skipDirs      bool
}

type IterOption func(*iterOptions)

func ContinueOnError(options *iterOptions) {
	options.continueOnErr = true
}

func SkipDirs(options *iterOptions) {
	options.skipDirs = true
}

func FilterErrors(filter ErrFilter) IterOption {
	return func(options *iterOptions) {
		options.errFilter = filter
	}
}

func Iter(fsys afero.Fs, root string, options ...IterOption) iter.Seq3[string, fs.FileInfo, error] {
	opts := iterOptions{}
	option.ApplyAll(&opts, options)

	return func(yield func(string, fs.FileInfo, error) bool) {
		done := false
		err := afero.Walk(fsys, root,
			func(path string, info fs.FileInfo, err error) error {
				if err != nil && !opts.continueOnErr {
					return err
				}
				if err != nil && opts.errFilter != nil {
					return opts.errFilter(err)
				}
				if info.IsDir() && opts.skipDirs {
					return nil
				}
				if done = !yield(path, info, err); done {
					return fs.SkipAll
				} else {
					return nil
				}
			},
		)
		if err != nil && !done {
			yield("", nil, err)
		}
	}
}
