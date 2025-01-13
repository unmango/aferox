package aferox

import (
	"io/fs"

	"github.com/spf13/afero"
)

type FoldFunc[T any] func(string, fs.FileInfo, T, error) (T, error)

func Fold[T any](fsys afero.Fs, root string, folder FoldFunc[T], initial T) (acc T, err error) {
	acc = initial
	err = afero.Walk(fsys, root,
		func(path string, info fs.FileInfo, err error) error {
			acc, err = folder(path, info, acc, err)
			return err
		},
	)

	return
}
