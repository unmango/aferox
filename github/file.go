package github

import (
	"errors"
	"io"
	"io/fs"

	"github.com/unmango/aferox/github/internal"
)

type File struct {
	internal.ReadOnlyFile
	name string
	file fs.File
}

// Close implements afero.File.
func (f *File) Close() error {
	return f.file.Close()
}

// Name implements afero.File.
func (f *File) Name() string {
	return f.name
}

// Read implements afero.File.
func (f *File) Read(p []byte) (int, error) {
	return f.file.Read(p)
}

// ReadAt implements afero.File.
func (f *File) ReadAt(p []byte, off int64) (int, error) {
	if ra, ok := f.file.(io.ReaderAt); ok {
		return ra.ReadAt(p, off)
	}
	return 0, errors.ErrUnsupported
}

// Readdir implements afero.File.
func (f *File) Readdir(count int) ([]fs.FileInfo, error) {
	rd, ok := f.file.(fs.ReadDirFile)
	if !ok {
		return nil, errors.ErrUnsupported
	}
	entries, err := rd.ReadDir(count)
	if err != nil {
		return nil, err
	}
	infos := make([]fs.FileInfo, len(entries))
	for i, e := range entries {
		info, err := e.Info()
		if err != nil {
			return nil, err
		}
		infos[i] = info
	}
	return infos, nil
}

// Readdirnames implements afero.File.
func (f *File) Readdirnames(n int) ([]string, error) {
	rd, ok := f.file.(fs.ReadDirFile)
	if !ok {
		return nil, errors.ErrUnsupported
	}
	entries, err := rd.ReadDir(n)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(entries))
	for i, e := range entries {
		names[i] = e.Name()
	}
	return names, nil
}

// Seek implements afero.File.
func (f *File) Seek(offset int64, whence int) (int64, error) {
	if s, ok := f.file.(io.Seeker); ok {
		return s.Seek(offset, whence)
	}
	return 0, errors.ErrUnsupported
}

// Stat implements afero.File.
func (f *File) Stat() (fs.FileInfo, error) {
	return f.file.Stat()
}
