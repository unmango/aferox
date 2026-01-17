package tar

import (
	"archive/tar"
	"os"
	"syscall"
	"time"

	"github.com/spf13/afero"
)

type Fs struct {
	w *tar.Writer
}

// Chmod implements [afero.Fs].
func (f *Fs) Chmod(name string, mode os.FileMode) error {
	return syscall.EPERM
}

// Chown implements [afero.Fs].
func (f *Fs) Chown(name string, uid int, gid int) error {
	return syscall.EPERM
}

// Chtimes implements [afero.Fs].
func (f *Fs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return syscall.EPERM
}

// Create implements [afero.Fs].
func (f *Fs) Create(name string) (afero.File, error) {
	return NewFile(name, f.w, 0), nil
}

// Mkdir implements [afero.Fs].
func (f *Fs) Mkdir(name string, perm os.FileMode) error {
	return nil
}

// MkdirAll implements [afero.Fs].
func (f *Fs) MkdirAll(path string, perm os.FileMode) error {
	return nil
}

// Name implements [afero.Fs].
func (f *Fs) Name() string {
	return "tar.Writer"
}

// Open implements [afero.Fs].
func (f *Fs) Open(name string) (afero.File, error) {
	return NewFile(name, f.w, 0), nil
}

// OpenFile implements [afero.Fs].
func (f *Fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	return NewFile(name, f.w, perm), nil
}

// Remove implements [afero.Fs].
func (f *Fs) Remove(name string) error {
	return syscall.EROFS
}

// RemoveAll implements [afero.Fs].
func (f *Fs) RemoveAll(path string) error {
	return syscall.EROFS
}

// Rename implements [afero.Fs].
func (f *Fs) Rename(oldname string, newname string) error {
	return syscall.EROFS
}

// Stat implements [afero.Fs].
func (f *Fs) Stat(name string) (os.FileInfo, error) {
	return &FileInfo{name: name, w: f.w}, nil
}

func NewFs(w *tar.Writer) afero.Fs {
	return &Fs{w: w}
}
