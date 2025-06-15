package mapped

import (
	"fmt"
	"maps"
	"os"
	"time"

	"github.com/spf13/afero"
)

type Fs map[string]afero.Fs

func NewFs(m map[string]afero.Fs) afero.Fs {
	fs := Fs{}
	for k, v := range m {
		fs[Clean(k)] = v
	}

	return fs
}

// Chmod implements afero.Fs.
func (f Fs) Chmod(name string, mode os.FileMode) error {
	if k, p, err := f.split(name); err != nil {
		return err
	} else {
		return f[k].Chmod(p, mode)
	}
}

// Chown implements afero.Fs.
func (f Fs) Chown(name string, uid int, gid int) error {
	if k, p, err := f.split(name); err != nil {
		return err
	} else {
		return f[k].Chown(p, uid, gid)
	}
}

// Chtimes implements afero.Fs.
func (f Fs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	if k, p, err := f.split(name); err != nil {
		return err
	} else {
		return f[k].Chtimes(p, atime, mtime)
	}
}

// Create implements afero.Fs.
func (f Fs) Create(name string) (afero.File, error) {
	if k, p, err := f.split(name); err != nil {
		return nil, err
	} else {
		return f[k].Create(p)
	}
}

// Mkdir implements afero.Fs.
func (f Fs) Mkdir(name string, perm os.FileMode) error {
	if k, p, err := f.split(name); err != nil {
		return err
	} else {
		return f[k].Mkdir(p, perm)
	}
}

// MkdirAll implements afero.Fs.
func (f Fs) MkdirAll(path string, perm os.FileMode) error {
	if k, p, err := f.split(path); err != nil {
		return err
	} else {
		return f[k].MkdirAll(p, perm)
	}
}

// Name implements afero.Fs.
func (f Fs) Name() string {
	return "Mapped"
}

// Open implements afero.Fs.
func (f Fs) Open(name string) (afero.File, error) {
	if k, p, err := f.split(name); err != nil {
		return nil, err
	} else {
		return f[k].Open(p)
	}
}

// OpenFile implements afero.Fs.
func (f Fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	if k, p, err := f.split(name); err != nil {
		return nil, err
	} else {
		return f[k].OpenFile(p, flag, perm)
	}
}

// Remove implements afero.Fs.
func (f Fs) Remove(name string) error {
	if k, p, err := f.split(name); err != nil {
		return err
	} else {
		return f[k].Remove(p)
	}
}

// RemoveAll implements afero.Fs.
func (f Fs) RemoveAll(path string) error {
	if k, p, err := f.split(path); err != nil {
		return err
	} else {
		return f[k].RemoveAll(p)
	}
}

// Rename implements afero.Fs.
func (f Fs) Rename(oldname string, newname string) error {
	panic("unimplemented")
}

// Stat implements afero.Fs.
func (f Fs) Stat(name string) (os.FileInfo, error) {
	if k, p, err := f.split(name); err != nil {
		return nil, err
	} else {
		return f[k].Stat(p)
	}
}

func (f Fs) split(name string) (key, path string, err error) {
	for k := range maps.Keys(f) {
		if p, ok := CutPrefix(name, k); ok {
			return k, p, nil
		}
	}

	return "", "", fmt.Errorf("%w: %s", os.ErrNotExist, name)
}
