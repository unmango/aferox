package filter

import (
	"io/fs"
	"path/filepath"

	"github.com/spf13/afero"
)

type File struct {
	file   afero.File
	filter Filter
}

// Close implements afero.File.
func (f *File) Close() error {
	return f.file.Close()
}

// Name implements afero.File.
func (f *File) Name() string {
	return f.file.Name()
}

// Read implements afero.File.
func (f *File) Read(p []byte) (n int, err error) {
	return f.file.Read(p)
}

// ReadAt implements afero.File.
func (f *File) ReadAt(p []byte, off int64) (n int, err error) {
	return f.file.ReadAt(p, off)
}

// Readdir implements afero.File.
func (f *File) Readdir(count int) (res []fs.FileInfo, err error) {
	var infos []fs.FileInfo
	infos, err = f.file.Readdir(count)
	if err != nil {
		return nil, err
	}

	for _, i := range infos {
		path := filepath.Join(f.file.Name(), i.Name())
		op := ReaddirOp{Name: path, Count: count}
		if i.IsDir() || f.filter(op) == nil {
			res = append(res, i)
		}
	}

	return res, nil
}

// Readdirnames implements afero.File.
func (f *File) Readdirnames(n int) (names []string, err error) {
	infos, err := f.Readdir(n)
	if err != nil {
		return nil, err
	}

	for _, i := range infos {
		names = append(names, i.Name())
	}

	return names, nil
}

// Seek implements afero.File.
func (f *File) Seek(offset int64, whence int) (int64, error) {
	return f.file.Seek(offset, whence)
}

// Stat implements afero.File.
func (f *File) Stat() (fs.FileInfo, error) {
	return f.file.Stat()
}

// Sync implements afero.File.
func (f *File) Sync() error {
	return f.file.Sync()
}

// Truncate implements afero.File.
func (f *File) Truncate(size int64) error {
	return f.file.Truncate(size)
}

// Write implements afero.File.
func (f *File) Write(p []byte) (n int, err error) {
	return f.file.Write(p)
}

// WriteAt implements afero.File.
func (f *File) WriteAt(p []byte, off int64) (n int, err error) {
	return f.file.WriteAt(p, off)
}

// WriteString implements afero.File.
func (f *File) WriteString(s string) (ret int, err error) {
	return f.file.WriteString(s)
}
