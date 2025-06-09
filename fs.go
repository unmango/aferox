package aferox

import (
	"io"

	"github.com/spf13/afero"
	"github.com/unmango/aferox/writer"
)

type (
	Writer = writer.Fs
)

func NewWriter(w io.Writer) afero.Fs {
	return writer.NewFs(w)
}
