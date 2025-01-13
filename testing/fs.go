package testing

import (
	"io/fs"
	"time"

	"github.com/spf13/afero"
)

type Fs struct {
	ChmodFunc     func(string, fs.FileMode) error
	ChownFunc     func(string, int, int) error
	ChtimesFunc   func(string, time.Time, time.Time) error
	CreateFunc    func(string) (afero.File, error)
	MkdirAllFunc  func(string, fs.FileMode) error
	MkdirFunc     func(string, fs.FileMode) error
	OpenFunc      func(string) (afero.File, error)
	OpenFileFunc  func(string, int, fs.FileMode) (afero.File, error)
	RemoveAllFunc func(string) error
	RemoveFunc    func(string) error
	RenameFunc    func(string, string) error
	StatFunc      func(string) (fs.FileInfo, error)
}

// Chmod implements afero.Fs.
func (f *Fs) Chmod(name string, mode fs.FileMode) error {
	if f.ChmodFunc == nil {
		panic("unimplemented")
	}

	return f.ChmodFunc(name, mode)
}

// Chown implements afero.Fs.
func (f *Fs) Chown(name string, uid int, gid int) error {
	if f.ChownFunc == nil {
		panic("unimplemented")
	}

	return f.ChownFunc(name, uid, gid)
}

// Chtimes implements afero.Fs.
func (f *Fs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	if f.ChtimesFunc == nil {
		panic("unimplemented")
	}

	return f.ChtimesFunc(name, atime, mtime)
}

// Create implements afero.Fs.
func (f *Fs) Create(name string) (afero.File, error) {
	if f.CreateFunc == nil {
		panic("unimplemented")
	}

	return f.CreateFunc(name)
}

// Mkdir implements afero.Fs.
func (f *Fs) Mkdir(name string, perm fs.FileMode) error {
	if f.MkdirFunc == nil {
		panic("unimplemented")
	}

	return f.MkdirFunc(name, perm)
}

// MkdirAll implements afero.Fs.
func (f *Fs) MkdirAll(path string, perm fs.FileMode) error {
	if f.MkdirAllFunc == nil {
		panic("unimplemented")
	}

	return f.MkdirAllFunc(path, perm)
}

// Name implements afero.Fs.
func (f *Fs) Name() string {
	return "Testing"
}

// Open implements afero.Fs.
func (f *Fs) Open(name string) (afero.File, error) {
	if f.OpenFunc == nil {
		panic("unimplemented")
	}

	return f.OpenFunc(name)
}

// OpenFile implements afero.Fs.
func (f *Fs) OpenFile(name string, flag int, perm fs.FileMode) (afero.File, error) {
	if f.OpenFileFunc == nil {
		panic("unimplemented")
	}

	return f.OpenFileFunc(name, flag, perm)
}

// Remove implements afero.Fs.
func (f *Fs) Remove(name string) error {
	if f.RemoveFunc == nil {
		panic("unimplemented")
	}

	return f.RemoveFunc(name)
}

// RemoveAll implements afero.Fs.
func (f *Fs) RemoveAll(path string) error {
	if f.RemoveAllFunc == nil {
		panic("unimplemented")
	}

	return f.RemoveAllFunc(path)
}

// Rename implements afero.Fs.
func (f *Fs) Rename(oldname string, newname string) error {
	if f.RenameFunc == nil {
		panic("unimplemented")
	}

	return f.RenameFunc(oldname, newname)
}

// Stat implements afero.Fs.
func (f *Fs) Stat(name string) (fs.FileInfo, error) {
	if f.StatFunc == nil {
		panic("unimplemented")
	}

	return f.StatFunc(name)
}

var _ afero.Fs = (*Fs)(nil)
