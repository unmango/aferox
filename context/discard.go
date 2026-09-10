package context

import (
	"io/fs"
	"time"

	"github.com/spf13/afero"
)

// DiscardFs adapts an [afero.Fs] to an [AferoFs] by ignoring the [context.Context] argument.
type DiscardFs struct{ afero.Fs }

// ChmodContext implements AferoFs.
func (a *DiscardFs) ChmodContext(ctx Context, name string, mode fs.FileMode) error {
	return a.Chmod(name, mode)
}

// ChownContext implements AferoFs.
func (a *DiscardFs) ChownContext(ctx Context, name string, uid int, gid int) error {
	return a.Chown(name, uid, gid)
}

// ChtimesContext implements AferoFs.
func (a *DiscardFs) ChtimesContext(ctx Context, name string, atime time.Time, mtime time.Time) error {
	return a.Chtimes(name, atime, mtime)
}

// CreateContext implements AferoFs.
func (a *DiscardFs) CreateContext(ctx Context, name string) (afero.File, error) {
	return a.Create(name)
}

// MkdirAllContext implements AferoFs.
func (a *DiscardFs) MkdirAllContext(ctx Context, path string, perm fs.FileMode) error {
	return a.MkdirAll(path, perm)
}

// MkdirContext implements AferoFs.
func (a *DiscardFs) MkdirContext(ctx Context, name string, perm fs.FileMode) error {
	return a.Mkdir(name, perm)
}

// OpenContext implements AferoFs.
func (a *DiscardFs) OpenContext(ctx Context, name string) (afero.File, error) {
	return a.Open(name)
}

// OpenFileContext implements AferoFs.
func (a *DiscardFs) OpenFileContext(ctx Context, name string, flag int, perm fs.FileMode) (afero.File, error) {
	return a.OpenFile(name, flag, perm)
}

// RemoveAllContext implements AferoFs.
func (a *DiscardFs) RemoveAllContext(ctx Context, path string) error {
	return a.RemoveAll(path)
}

// RemoveContext implements AferoFs.
func (a *DiscardFs) RemoveContext(ctx Context, name string) error {
	return a.Remove(name)
}

// RenameContext implements AferoFs.
func (a *DiscardFs) RenameContext(ctx Context, oldname string, newname string) error {
	return a.Rename(oldname, newname)
}

// StatContext implements AferoFs.
func (a *DiscardFs) StatContext(ctx Context, name string) (fs.FileInfo, error) {
	return a.Stat(name)
}
