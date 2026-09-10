package testing

import (
	"io/fs"
	"time"

	"github.com/spf13/afero"
	"github.com/unmango/aferox/context"
)

type ContextFs struct {
	ChmodFunc     func(context.Context, string, fs.FileMode) error
	ChownFunc     func(context.Context, string, int, int) error
	ChtimesFunc   func(context.Context, string, time.Time, time.Time) error
	CreateFunc    func(context.Context, string) (afero.File, error)
	MkdirFunc     func(context.Context, string, fs.FileMode) error
	MkdirAllFunc  func(context.Context, string, fs.FileMode) error
	OpenFunc      func(context.Context, string) (afero.File, error)
	OpenFileFunc  func(context.Context, string, int, fs.FileMode) (afero.File, error)
	RemoveFunc    func(context.Context, string) error
	RemoveAllFunc func(context.Context, string) error
	RenameFunc    func(context.Context, string, string) error
	StatFunc      func(context.Context, string) (fs.FileInfo, error)
}

// Chmod implements context.Fs.
func (c *ContextFs) Chmod(ctx context.Context, name string, mode fs.FileMode) error {
	if c.ChmodFunc == nil {
		panic("unimplemented")
	}

	return c.ChmodFunc(ctx, name, mode)
}

// Chown implements context.Fs.
func (c *ContextFs) Chown(ctx context.Context, name string, uid int, gid int) error {
	if c.ChownFunc == nil {
		panic("unimplemented")
	}

	return c.ChownFunc(ctx, name, uid, gid)
}

// Chtimes implements context.Fs.
func (c *ContextFs) Chtimes(ctx context.Context, name string, atime time.Time, mtime time.Time) error {
	if c.ChtimesFunc == nil {
		panic("unimplemented")
	}

	return c.ChtimesFunc(ctx, name, atime, mtime)
}

// Create implements context.Fs.
func (c *ContextFs) Create(ctx context.Context, name string) (afero.File, error) {
	if c.CreateFunc == nil {
		panic("unimplemented")
	}

	return c.CreateFunc(ctx, name)
}

// Mkdir implements context.Fs.
func (c *ContextFs) Mkdir(ctx context.Context, name string, perm fs.FileMode) error {
	if c.MkdirFunc == nil {
		panic("unimplemented")
	}

	return c.MkdirFunc(ctx, name, perm)
}

// MkdirAll implements context.Fs.
func (c *ContextFs) MkdirAll(ctx context.Context, path string, perm fs.FileMode) error {
	if c.MkdirAllFunc == nil {
		panic("unimplemented")
	}

	return c.MkdirAllFunc(ctx, path, perm)
}

// Name implements context.Fs.
func (c *ContextFs) Name() string {
	return "Testing"
}

// Open implements context.Fs.
func (c *ContextFs) Open(ctx context.Context, name string) (afero.File, error) {
	if c.OpenFunc == nil {
		panic("unimplemented")
	}

	return c.OpenFunc(ctx, name)
}

// OpenFile implements context.Fs.
func (c *ContextFs) OpenFile(ctx context.Context, name string, flag int, perm fs.FileMode) (afero.File, error) {
	if c.OpenFileFunc == nil {
		panic("unimplemented")
	}

	return c.OpenFileFunc(ctx, name, flag, perm)
}

// Remove implements context.Fs.
func (c *ContextFs) Remove(ctx context.Context, name string) error {
	if c.RemoveFunc == nil {
		panic("unimplemented")
	}

	return c.RemoveFunc(ctx, name)
}

// RemoveAll implements context.Fs.
func (c *ContextFs) RemoveAll(ctx context.Context, path string) error {
	if c.RemoveAllFunc == nil {
		panic("unimplemented")
	}

	return c.RemoveAllFunc(ctx, path)
}

// Rename implements context.Fs.
func (c *ContextFs) Rename(ctx context.Context, oldname string, newname string) error {
	if c.RenameFunc == nil {
		panic("unimplemented")
	}

	return c.RenameFunc(ctx, oldname, newname)
}

// Stat implements context.Fs.
func (c *ContextFs) Stat(ctx context.Context, name string) (fs.FileInfo, error) {
	if c.StatFunc == nil {
		panic("unimplemented")
	}

	return c.StatFunc(ctx, name)
}

var _ context.Fs = (*ContextFs)(nil)
