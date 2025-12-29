package containers

import (
	"os"
	"time"

	"github.com/spf13/afero"
	"go.podman.io/storage"
)

type Fs struct {
	store storage.Store
}

// Chmod implements [afero.Fs].
func (f *Fs) Chmod(name string, mode os.FileMode) error {
	panic("unimplemented")
}

// Chown implements [afero.Fs].
func (f *Fs) Chown(name string, uid int, gid int) error {
	panic("unimplemented")
}

// Chtimes implements [afero.Fs].
func (f *Fs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	panic("unimplemented")
}

// Create implements [afero.Fs].
func (f *Fs) Create(name string) (afero.File, error) {
	panic("unimplemented")
}

// Mkdir implements [afero.Fs].
func (f *Fs) Mkdir(name string, perm os.FileMode) error {
	panic("unimplemented")
}

// MkdirAll implements [afero.Fs].
func (f *Fs) MkdirAll(path string, perm os.FileMode) error {
	panic("unimplemented")
}

// Name implements [afero.Fs].
func (f *Fs) Name() string {
	panic("unimplemented")
}

// Open implements [afero.Fs].
func (f *Fs) Open(name string) (afero.File, error) {
	panic("unimplemented")
}

// OpenFile implements [afero.Fs].
func (f *Fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	panic("unimplemented")
}

// Remove implements [afero.Fs].
func (f *Fs) Remove(name string) error {
	panic("unimplemented")
}

// RemoveAll implements [afero.Fs].
func (f *Fs) RemoveAll(path string) error {
	panic("unimplemented")
}

// Rename implements [afero.Fs].
func (f *Fs) Rename(oldname string, newname string) error {
	panic("unimplemented")
}

// Stat implements [afero.Fs].
func (f *Fs) Stat(name string) (os.FileInfo, error) {
	panic("unimplemented")
}

func NewFs(store storage.Store) afero.Fs {
	return &Fs{
		store: store,
	}
}
