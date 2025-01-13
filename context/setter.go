package context

import (
	"fmt"
	"io/fs"
	"time"

	"github.com/spf13/afero"
)

// WithSetterFs adapts an [afero.Fs] to an [Fs] by calling [SetContext] on the given
// [ContextSetter] before forwarding each operation to the given [afero.Fs]
type WithSetterFs struct {
	Setter
	Fs afero.Fs
}

// Chmod implements Fs.
func (s *WithSetterFs) Chmod(ctx Context, name string, mode fs.FileMode) error {
	s.SetContext(ctx)
	return s.Fs.Chmod(name, mode)
}

// Chown implements Fs.
func (s *WithSetterFs) Chown(ctx Context, name string, uid int, gid int) error {
	s.SetContext(ctx)
	return s.Fs.Chown(name, uid, gid)
}

// Chtimes implements Fs.
func (s *WithSetterFs) Chtimes(ctx Context, name string, atime time.Time, mtime time.Time) error {
	s.SetContext(ctx)
	return s.Fs.Chtimes(name, atime, mtime)
}

// Create implements Fs.
func (s *WithSetterFs) Create(ctx Context, name string) (afero.File, error) {
	s.SetContext(ctx)
	return s.Fs.Create(name)
}

// Mkdir implements Fs.
func (s *WithSetterFs) Mkdir(ctx Context, name string, perm fs.FileMode) error {
	s.SetContext(ctx)
	return s.Fs.Mkdir(name, perm)
}

// MkdirAll implements Fs.
func (s *WithSetterFs) MkdirAll(ctx Context, path string, perm fs.FileMode) error {
	s.SetContext(ctx)
	return s.Fs.MkdirAll(path, perm)
}

// Name implements Fs.
func (s *WithSetterFs) Name() string {
	return fmt.Sprintf("Scoped: %s", s.Fs.Name())
}

// Open implements Fs.
func (s *WithSetterFs) Open(ctx Context, name string) (afero.File, error) {
	s.SetContext(ctx)
	return s.Fs.Open(name)
}

// OpenFile implements Fs.
func (s *WithSetterFs) OpenFile(ctx Context, name string, flag int, perm fs.FileMode) (afero.File, error) {
	s.SetContext(ctx)
	return s.Fs.OpenFile(name, flag, perm)
}

// Remove implements Fs.
func (s *WithSetterFs) Remove(ctx Context, name string) error {
	s.SetContext(ctx)
	return s.Fs.Remove(name)
}

// RemoveAll implements Fs.
func (s *WithSetterFs) RemoveAll(ctx Context, path string) error {
	s.SetContext(ctx)
	return s.Fs.RemoveAll(path)
}

// Rename implements Fs.
func (s *WithSetterFs) Rename(ctx Context, oldname string, newname string) error {
	s.SetContext(ctx)
	return s.Fs.Rename(oldname, newname)
}

// Stat implements Fs.
func (s *WithSetterFs) Stat(ctx Context, name string) (fs.FileInfo, error) {
	s.SetContext(ctx)
	return s.Fs.Stat(name)
}

var _ Fs = (*WithSetterFs)(nil)
