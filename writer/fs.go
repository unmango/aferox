package writer

import (
	"io"
	"io/fs"

	"github.com/spf13/afero"
)

type Fs struct {
	afero.ReadOnlyFs
	writer io.Writer
}

// Name implements afero.Fs.
func (w *Fs) Name() string {
	return "io.Writer"
}

// Open implements afero.Fs.
func (w *Fs) Open(name string) (afero.File, error) {
	return &File{w.writer, name}, nil
}

// OpenFile implements afero.Fs.
func (w *Fs) OpenFile(name string, _ int, _ fs.FileMode) (afero.File, error) {
	return &File{w.writer, name}, nil
}

// Stat implements afero.Fs.
func (w *Fs) Stat(name string) (fs.FileInfo, error) {
	return &FileInfo{w.writer, name}, nil
}

func NewFs(writer io.Writer) afero.Fs {
	return &Fs{writer: writer}
}
