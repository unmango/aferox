package filter

import (
	"fmt"
	"io/fs"
	"syscall"
	"time"

	"github.com/spf13/afero"
	"github.com/unmango/aferox/op"
)

type (
	// Filter evaluates a filesystem operation and returns an error if it should be blocked.
	Filter func(op.Operation) error

	// Predicate evaluates a filesystem operation and returns true if it should be allowed.
	Predicate func(op.Operation) bool
)

func FromFilter(base afero.Fs, filter Filter) afero.Fs {
	return &Fs{src: base, filter: filter}
}

func FromPredicateWithError(base afero.Fs, pred Predicate, onFalse error) afero.Fs {
	return FromFilter(base, func(operation op.Operation) error {
		if pred(operation) {
			return nil
		}
		return onFalse
	})
}

func FromPredicate(base afero.Fs, pred Predicate) afero.Fs {
	return FromPredicateWithError(base, pred, syscall.ENOENT)
}

func NewFs(base afero.Fs, predicate Predicate) afero.Fs {
	return FromPredicate(base, predicate)
}

// PathFilter creates a Filter from a path-based filter function.
// This is a convenience function for backward compatibility.
func PathFilter(fn func(string) error) Filter {
	return func(operation op.Operation) error {
		return fn(operation.Path())
	}
}

// PathPredicate creates a Predicate from a path-based predicate function.
// This is a convenience function for backward compatibility.
func PathPredicate(fn func(string) bool) Predicate {
	return func(operation op.Operation) bool {
		return fn(operation.Path())
	}
}

type Fs struct {
	src    afero.Fs
	filter Filter
}

// Chmod implements afero.Fs.
func (f *Fs) Chmod(name string, mode fs.FileMode) error {
	if err := f.dirOrMatches(op.Chmod{
		Name: name,
		Mode: mode,
	}); err != nil {
		return err
	}

	return f.src.Chmod(name, mode)
}

// Chown implements afero.Fs.
func (f *Fs) Chown(name string, uid int, gid int) error {
	if err := f.dirOrMatches(op.Chown{
		Name: name,
		UID:  uid,
		GID:  gid,
	}); err != nil {
		return err
	}

	return f.src.Chown(name, uid, gid)
}

// Chtimes implements afero.Fs.
func (f *Fs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	if err := f.dirOrMatches(op.Chtimes{
		Name:  name,
		Atime: atime,
		Mtime: mtime,
	}); err != nil {
		return err
	}
	return f.src.Chtimes(name, atime, mtime)
}

// Create implements afero.Fs.
func (f *Fs) Create(name string) (afero.File, error) {
	if err := f.matches(op.Create{Name: name}); err != nil {
		return nil, err
	}
	return f.src.Create(name)
}

// Mkdir implements afero.Fs.
func (f *Fs) Mkdir(name string, perm fs.FileMode) error {
	return f.src.Mkdir(name, perm)
}

// MkdirAll implements afero.Fs.
func (f *Fs) MkdirAll(path string, perm fs.FileMode) error {
	return f.src.Mkdir(path, perm)
}

// Name implements afero.Fs.
func (f *Fs) Name() string {
	return fmt.Sprintf("Filter: %s", f.src.Name())
}

// Open implements afero.Fs.
func (f *Fs) Open(name string) (afero.File, error) {
	dir, err := afero.IsDir(f.src, name)
	if err != nil {
		return nil, err
	}
	if !dir {
		if err := f.matches(OpenOp{Name: name}); err != nil {
			return nil, err
		}
	}

	file, err := f.src.Open(name)
	if err != nil {
		return nil, err
	}

	return &File{
		file:   file,
		filter: f.filter,
	}, nil
}

// OpenFile implements afero.Fs.
func (f *Fs) OpenFile(name string, flag int, perm fs.FileMode) (afero.File, error) {
	if err := f.dirOrMatches(op.OpenFile{
		Name: name,
		Flag: flag,
		Perm: perm,
	}); err != nil {
		return nil, err
	}
	return f.src.OpenFile(name, flag, perm)
}

// Remove implements afero.Fs.
func (f *Fs) Remove(name string) error {
	if err := f.dirOrMatches(op.Remove{Name: name}); err != nil {
		return err
	}
	return f.src.Remove(name)
}

// RemoveAll implements afero.Fs.
func (f *Fs) RemoveAll(path string) error {
	dir, err := afero.IsDir(f.src, path)
	if err != nil {
		return err
	}
	if !dir {
		if err = f.matches(RemoveAllOp{PathName: path}); err != nil {
			return err
		}
	}

	return f.src.RemoveAll(path)
}

// Rename implements afero.Fs.
func (f *Fs) Rename(oldname string, newname string) error {
	dir, err := afero.IsDir(f.src, oldname)
	if err != nil {
		return err
	}
	if dir {
		return nil
	}
	op := RenameOp{Oldname: oldname, Newname: newname}
	if err = f.matches(op); err != nil {
		return err
	}
	// Also check the new name
	if err = f.matches(CreateOp{Name: newname}); err != nil {
		return err
	}

	return f.src.Rename(oldname, newname)
}

// Stat implements afero.Fs.
func (f *Fs) Stat(name string) (fs.FileInfo, error) {
	if err := f.dirOrMatches(StatOp{Name: name}); err != nil {
		return nil, err
	}

	return f.src.Stat(name)
}

func (f *Fs) dirOrMatches(op Operation) error {
	dir, err := afero.IsDir(f.src, op.Path())
	if err != nil {
		return err
	}
	if dir {
		return nil
	}

	return f.matches(operation)
}

func (f *Fs) matches(operation op.Operation) error {
	if f.filter == nil {
		return nil
	} else {
		return f.filter(operation)
	}
}
