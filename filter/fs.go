package filter

import (
	"fmt"
	"io/fs"
	"syscall"
	"time"

	"github.com/spf13/afero"
)

type (
	Filter    func(string) error
	Predicate func(string) bool
)

func FromFilter(base afero.Fs, filter Filter) afero.Fs {
	return &Fs{src: base, filter: filter}
}

func FromPredicateWithError(base afero.Fs, pred Predicate, onFalse error) afero.Fs {
	return FromFilter(base, func(s string) error {
		if pred(s) {
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

type Fs struct {
	src    afero.Fs
	filter Filter
}

// Chmod implements afero.Fs.
func (f *Fs) Chmod(name string, mode fs.FileMode) error {
	if err := f.dirOrMatches(name); err != nil {
		return err
	}

	return f.src.Chmod(name, mode)
}

// Chown implements afero.Fs.
func (f *Fs) Chown(name string, uid int, gid int) error {
	if err := f.dirOrMatches(name); err != nil {
		return err
	}

	return f.src.Chown(name, uid, gid)
}

// Chtimes implements afero.Fs.
func (f *Fs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	if err := f.dirOrMatches(name); err != nil {
		return err
	}

	return f.src.Chtimes(name, atime, mtime)
}

// Create implements afero.Fs.
func (f *Fs) Create(name string) (afero.File, error) {
	if err := f.matchesName(name); err != nil {
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
		if err := f.matchesName(name); err != nil {
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
	if err := f.dirOrMatches(name); err != nil {
		return nil, err
	}

	return f.src.OpenFile(name, flag, perm)
}

// Remove implements afero.Fs.
func (f *Fs) Remove(name string) error {
	if err := f.dirOrMatches(name); err != nil {
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
		if err = f.matchesName(path); err != nil {
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
	if err = f.matchesName(oldname); err != nil {
		return err
	}
	if err = f.matchesName(newname); err != nil {
		return err
	}

	return f.src.Rename(oldname, newname)
}

// Stat implements afero.Fs.
func (f *Fs) Stat(name string) (fs.FileInfo, error) {
	if err := f.dirOrMatches(name); err != nil {
		return nil, err
	}

	return f.src.Stat(name)
}

func (f *Fs) dirOrMatches(name string) error {
	dir, err := afero.IsDir(f.src, name)
	if err != nil {
		return err
	}
	if dir {
		return nil
	}

	return f.matchesName(name)
}

func (f *Fs) matchesName(name string) error {
	if f.filter == nil {
		return nil
	} else {
		return f.filter(name)
	}
}
