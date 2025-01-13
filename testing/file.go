package testing

import (
	"io/fs"

	"github.com/spf13/afero"
)

type File struct {
	CloseFunc        func() error
	NameFunc         func() string
	ReadFunc         func([]byte) (int, error)
	ReadAtFunc       func([]byte, int64) (int, error)
	ReaddirFunc      func(int) ([]fs.FileInfo, error)
	ReaddirnamesFunc func(int) ([]string, error)
	SeekFunc         func(int64, int) (int64, error)
	StatFunc         func() (fs.FileInfo, error)
	SyncFunc         func() error
	TruncateFunc     func(int64) error
	WriteFunc        func([]byte) (int, error)
	WriteAtFunc      func([]byte, int64) (int, error)
	WriteStringFunc  func(string) (int, error)
}

// Close implements afero.File.
func (f *File) Close() error {
	if f.CloseFunc == nil {
		panic("unimplemented")
	}

	return f.CloseFunc()
}

// Name implements afero.File.
func (f *File) Name() string {
	if f.NameFunc == nil {
		panic("unimplemented")
	}

	return f.NameFunc()
}

// Read implements afero.File.
func (f *File) Read(p []byte) (n int, err error) {
	if f.ReadFunc == nil {
		panic("unimplemented")
	}

	return f.ReadFunc(p)
}

// ReadAt implements afero.File.
func (f *File) ReadAt(p []byte, off int64) (n int, err error) {
	if f.ReadAtFunc == nil {
		panic("unimplemented")
	}

	return f.ReadAtFunc(p, off)
}

// Readdir implements afero.File.
func (f *File) Readdir(count int) ([]fs.FileInfo, error) {
	if f.ReaddirFunc == nil {
		panic("unimplemented")
	}

	return f.ReaddirFunc(count)
}

// Readdirnames implements afero.File.
func (f *File) Readdirnames(n int) ([]string, error) {
	if f.ReaddirnamesFunc == nil {
		panic("unimplemented")
	}

	return f.ReaddirnamesFunc(n)
}

// Seek implements afero.File.
func (f *File) Seek(offset int64, whence int) (int64, error) {
	if f.SeekFunc == nil {
		panic("unimplemented")
	}

	return f.SeekFunc(offset, whence)
}

// Stat implements afero.File.
func (f *File) Stat() (fs.FileInfo, error) {
	if f.StatFunc == nil {
		panic("unimplemented")
	}

	return f.StatFunc()
}

// Sync implements afero.File.
func (f *File) Sync() error {
	if f.SyncFunc == nil {
		panic("unimplemented")
	}

	return f.SyncFunc()
}

// Truncate implements afero.File.
func (f *File) Truncate(size int64) error {
	if f.TruncateFunc == nil {
		panic("unimplemented")
	}

	return f.TruncateFunc(size)
}

// Write implements afero.File.
func (f *File) Write(p []byte) (n int, err error) {
	if f.WriteFunc == nil {
		panic("unimplemented")
	}

	return f.WriteFunc(p)
}

// WriteAt implements afero.File.
func (f *File) WriteAt(p []byte, off int64) (n int, err error) {
	if f.WriteAtFunc == nil {
		panic("unimplemented")
	}

	return f.WriteAtFunc(p, off)
}

// WriteString implements afero.File.
func (f *File) WriteString(s string) (ret int, err error) {
	if f.WriteStringFunc == nil {
		panic("unimplemented")
	}

	return f.WriteStringFunc(s)
}

var _ afero.File = (*File)(nil)
