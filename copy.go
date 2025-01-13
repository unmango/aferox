package aferox

import (
	"github.com/spf13/afero"
	"github.com/unmango/aferox/internal"
)

func Copy(src, dest afero.Fs) error {
	return internal.Copy(src, dest)
}
