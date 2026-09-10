package internal

import (
	"fmt"
	"io/fs"

	"github.com/spf13/afero"
)

func Copy(src, dest afero.Fs) error {
	return afero.Walk(src, "",
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if path == "" {
				return nil // Skip root
			}
			if info.IsDir() {
				return dest.Mkdir(path, info.Mode())
			}

			if f, err := src.Open(path); err != nil {
				return fmt.Errorf("open %s: %w", path, err)
			} else {
				return afero.WriteReader(dest, path, f)
			}
		},
	)
}
