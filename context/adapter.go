package context

import (
	"io/fs"
	"time"

	"github.com/spf13/afero"
)

// AdapterFs adapts an [Fs] to a partial [AferoFs] structure.
// It is intended to be embedded as a utility type to satisfy
// an [AferoFs] interface with an [Fs]
type AdapterFs struct{ Fs Fs }

// ChmodContext implements AferoFs.
func (a *AdapterFs) ChmodContext(ctx Context, name string, mode fs.FileMode) error {
	return a.Fs.Chmod(ctx, name, mode)
}

// ChownContext implements AferoFs.
func (a *AdapterFs) ChownContext(ctx Context, name string, uid int, gid int) error {
	return a.Fs.Chown(ctx, name, uid, gid)
}

// ChtimesContext implements AferoFs.
func (a *AdapterFs) ChtimesContext(ctx Context, name string, atime time.Time, mtime time.Time) error {
	return a.Fs.Chtimes(ctx, name, atime, mtime)
}

// CreateContext implements AferoFs.
func (a *AdapterFs) CreateContext(ctx Context, name string) (afero.File, error) {
	return a.Fs.Create(ctx, name)
}

// MkdirAllContext implements AferoFs.
func (a *AdapterFs) MkdirAllContext(ctx Context, path string, perm fs.FileMode) error {
	return a.Fs.MkdirAll(ctx, path, perm)
}

// MkdirContext implements AferoFs.
func (a *AdapterFs) MkdirContext(ctx Context, name string, perm fs.FileMode) error {
	return a.Fs.Mkdir(ctx, name, perm)
}

// OpenContext implements AferoFs.
func (a *AdapterFs) OpenContext(ctx Context, name string) (afero.File, error) {
	return a.Fs.Open(ctx, name)
}

// OpenFileContext implements AferoFs.
func (a *AdapterFs) OpenFileContext(ctx Context, name string, flag int, perm fs.FileMode) (afero.File, error) {
	return a.Fs.OpenFile(ctx, name, flag, perm)
}

// RemoveAllContext implements AferoFs.
func (a *AdapterFs) RemoveAllContext(ctx Context, path string) error {
	return a.Fs.RemoveAll(ctx, path)
}

// RemoveContext implements AferoFs.
func (a *AdapterFs) RemoveContext(ctx Context, name string) error {
	return a.Fs.Remove(ctx, name)
}

// RenameContext implements AferoFs.
func (a *AdapterFs) RenameContext(ctx Context, oldname string, newname string) error {
	return a.Fs.Rename(ctx, oldname, newname)
}

// StatContext implements AferoFs.
func (a *AdapterFs) StatContext(ctx Context, name string) (fs.FileInfo, error) {
	return a.Fs.Stat(ctx, name)
}
