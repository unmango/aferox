package aferox

import (
	"errors"
	"io/fs"

	"github.com/spf13/afero"
	"github.com/unmango/go/option"
)

func StatFirst(fsys afero.Fs, root string, options ...IterOption) (fs.FileInfo, error) {
	opts := &iterOptions{}
	option.ApplyAll(opts, options)

	info, err := Fold(fsys, root,
		func(path string, info fs.FileInfo, acc fs.FileInfo, err error) (fs.FileInfo, error) {
			if err != nil {
				return nil, err
			}
			if path == "." || path == "" {
				return nil, nil
			}
			if info.IsDir() && opts.skipDirs {
				return nil, nil
			}

			return info, fs.SkipAll
		},
		nil,
	)

	if err != nil && !errors.Is(err, fs.SkipAll) {
		return nil, err
	}
	if info == nil {
		return nil, errors.New("Fs contains no entries")
	}

	return info, nil
}

func OpenFirst(fsys afero.Fs, root string, options ...IterOption) (afero.File, error) {
	opts := &iterOptions{}
	option.ApplyAll(opts, options)

	file, err := Fold(fsys, root,
		func(path string, info fs.FileInfo, acc afero.File, err error) (afero.File, error) {
			if err != nil {
				return nil, err
			}
			if path == "." || path == "" {
				return nil, nil
			}
			if info.IsDir() && opts.skipDirs {
				return nil, nil
			}

			if file, err := fsys.Open(path); err != nil {
				return nil, err
			} else {
				return file, fs.SkipAll
			}
		},
		nil,
	)

	if err != nil && !errors.Is(err, fs.SkipAll) {
		return nil, err
	}
	if file == nil {
		return nil, errors.New("Fs contains no entries")
	}

	return file, nil
}
