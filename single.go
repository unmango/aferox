package aferox

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/spf13/afero"
	"github.com/unmango/go/option"
)

type errSingle struct {
	acc string
	cur string
}

func (e errSingle) Error() string {
	return fmt.Sprintf("fs contains more than one entry\n\thad:\t%s\n\tfound:\t%s", e.acc, e.cur)
}

func StatSingle(fsys afero.Fs, root string, options ...IterOption) (fs.FileInfo, error) {
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
			if acc != nil {
				return nil, errSingle{acc.Name(), info.Name()}
			}

			return info, nil
		},
		nil,
	)

	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, errors.New("Fs contains no entries")
	}

	return info, nil
}

func OpenSingle(fsys afero.Fs, root string, options ...IterOption) (afero.File, error) {
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
			if acc != nil {
				return nil, errSingle{acc.Name(), info.Name()}
			}

			return fsys.Open(path)
		},
		nil,
	)

	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, errors.New("Fs contains no entries")
	}

	return file, nil
}
