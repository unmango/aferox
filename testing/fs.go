package testing

import (
	"io/fs"
	"time"

	"github.com/spf13/afero"
)

type Fs struct {
	afero.Fs

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
	if f.ChmodFunc != nil {
		return f.ChmodFunc(name, mode)
	}
	return f.fs().Chmod(name, mode)
}

// Chown implements afero.Fs.
func (f *Fs) Chown(name string, uid int, gid int) error {
	if f.ChownFunc != nil {
		return f.ChownFunc(name, uid, gid)
	}
	return f.fs().Chown(name, uid, gid)
}

// Chtimes implements afero.Fs.
func (f *Fs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	if f.ChtimesFunc != nil {
		return f.ChtimesFunc(name, atime, mtime)
	}
	return f.fs().Chtimes(name, atime, mtime)
}

// Create implements afero.Fs.
func (f *Fs) Create(name string) (afero.File, error) {
	if f.CreateFunc != nil {
		return f.CreateFunc(name)
	}
	return f.fs().Create(name)
}

// Mkdir implements afero.Fs.
func (f *Fs) Mkdir(name string, perm fs.FileMode) error {
	if f.MkdirFunc != nil {
		return f.MkdirFunc(name, perm)
	}
	return f.fs().Mkdir(name, perm)
}

// MkdirAll implements afero.Fs.
func (f *Fs) MkdirAll(path string, perm fs.FileMode) error {
	if f.MkdirAllFunc != nil {
		return f.MkdirAllFunc(path, perm)
	}
	return f.fs().MkdirAll(path, perm)
}

// Name implements afero.Fs.
func (f *Fs) Name() string {
	return "Testing"
}

// Open implements afero.Fs.
func (f *Fs) Open(name string) (afero.File, error) {
	if f.OpenFunc != nil {
		return f.OpenFunc(name)
	}
	return f.fs().Open(name)
}

// OpenFile implements afero.Fs.
func (f *Fs) OpenFile(name string, flag int, perm fs.FileMode) (afero.File, error) {
	if f.OpenFileFunc != nil {
		return f.OpenFileFunc(name, flag, perm)
	}
	return f.fs().OpenFile(name, flag, perm)
}

// Remove implements afero.Fs.
func (f *Fs) Remove(name string) error {
	if f.RemoveFunc != nil {
		return f.RemoveFunc(name)
	}
	return f.fs().Remove(name)
}

// RemoveAll implements afero.Fs.
func (f *Fs) RemoveAll(path string) error {
	if f.RemoveAllFunc != nil {
		return f.RemoveAllFunc(path)
	}
	return f.fs().RemoveAll(path)
}

// Rename implements afero.Fs.
func (f *Fs) Rename(oldname string, newname string) error {
	if f.RenameFunc != nil {
		return f.RenameFunc(oldname, newname)
	}
	return f.fs().Rename(oldname, newname)
}

// Stat implements afero.Fs.
func (f *Fs) Stat(name string) (fs.FileInfo, error) {
	if f.StatFunc != nil {
		return f.StatFunc(name)
	}
	return f.fs().Stat(name)
}

func (f *Fs) fs() afero.Fs {
	if f.Fs == nil {
		f.Fs = afero.NewMemMapFs()
	}
	return f.Fs
}
