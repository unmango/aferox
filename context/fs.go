package context

import (
	"os"
	"time"

	"github.com/spf13/afero"
)

type File = afero.File

// Fs is a filesystem interface.
type Fs interface {
	// Create creates a file in the filesystem, returning the file and an
	// error, if any happens.
	Create(ctx Context, name string) (File, error)

	// Mkdir creates a directory in the filesystem, return an error if any
	// happens.
	Mkdir(ctx Context, name string, perm os.FileMode) error

	// MkdirAll creates a directory path and all parents that does not exist
	// yet.
	MkdirAll(ctx Context, path string, perm os.FileMode) error

	// Open opens a file, returning it or an error, if any happens.
	Open(ctx Context, name string) (File, error)

	// OpenFile opens a file using the given flags and the given mode.
	OpenFile(ctx Context, name string, flag int, perm os.FileMode) (File, error)

	// Remove removes a file identified by name, returning an error, if any
	// happens.
	Remove(ctx Context, name string) error

	// RemoveAll removes a directory path and any children it contains. It
	// does not fail if the path does not exist (return nil).
	RemoveAll(ctx Context, path string) error

	// Rename renames a file.
	Rename(ctx Context, oldname, newname string) error

	// Stat returns a FileInfo describing the named file, or an error, if any
	// happens.
	Stat(ctx Context, name string) (os.FileInfo, error)

	// The name of this FileSystem
	Name() string

	// Chmod changes the mode of the named file to mode.
	Chmod(ctx Context, name string, mode os.FileMode) error

	// Chown changes the uid and gid of the named file.
	Chown(ctx Context, name string, uid, gid int) error

	// Chtimes changes the access and modification times of the named file
	Chtimes(ctx Context, name string, atime time.Time, mtime time.Time) error
}

// AferoFs is a filesystem interface.
type AferoFs interface {
	afero.Fs

	// CreateContext creates a file in the filesystem, returning the file and an
	// error, if any happens.
	CreateContext(ctx Context, name string) (File, error)

	// MkdirContext creates a directory in the filesystem, return an error if any
	// happens.
	MkdirContext(ctx Context, name string, perm os.FileMode) error

	// MkdirAllContext creates a directory path and all parents that does not exist
	// yet.
	MkdirAllContext(ctx Context, path string, perm os.FileMode) error

	// OpenContext opens a file, returning it or an error, if any happens.
	OpenContext(ctx Context, name string) (File, error)

	// OpenFileContext opens a file using the given flags and the given mode.
	OpenFileContext(ctx Context, name string, flag int, perm os.FileMode) (File, error)

	// RemoveContext removes a file identified by name, returning an error, if any
	// happens.
	RemoveContext(ctx Context, name string) error

	// RemoveAllContext removes a directory path and any children it contains. It
	// does not fail if the path does not exist (return nil).
	RemoveAllContext(ctx Context, path string) error

	// RenameContext renames a file.
	RenameContext(ctx Context, oldname, newname string) error

	// StatContext returns a FileInfo describing the named file, or an error, if any
	// happens.
	StatContext(ctx Context, name string) (os.FileInfo, error)

	// ChmodContext changes the mode of the named file to mode.
	ChmodContext(ctx Context, name string, mode os.FileMode) error

	// ChownContext changes the uid and gid of the named file.
	ChownContext(ctx Context, name string, uid, gid int) error

	// ChtimesContext changes the access and modification times of the named file
	ChtimesContext(ctx Context, name string, atime time.Time, mtime time.Time) error
}
