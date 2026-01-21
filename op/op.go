package op

import (
	"io/fs"
	"time"
)

// Operation represents a filesystem operation being performed.
// Use a type switch to determine the specific operation type.
type Operation interface {
	Path() string
}

// Chmod represents a chmod operation.
type Chmod struct {
	Name string
	Mode fs.FileMode
}

func (o Chmod) Path() string { return o.Name }

// Chown represents a chown operation.
type Chown struct {
	Name string
	UID  int
	GID  int
}

func (o Chown) Path() string { return o.Name }

// Chtimes represents a chtimes operation.
type Chtimes struct {
	Name  string
	Atime time.Time
	Mtime time.Time
}

func (o Chtimes) Path() string { return o.Name }

// Create represents a create operation.
type Create struct {
	Name string
}

func (o Create) Path() string { return o.Name }

// Mkdir represents a mkdir operation.
type Mkdir struct {
	Name string
	Perm fs.FileMode
}

func (o Mkdir) Path() string { return o.Name }

// MkdirAll represents a mkdir -p operation.
type MkdirAll struct {
	PathName string
	Perm     fs.FileMode
}

func (o MkdirAll) Path() string { return o.PathName }

// Open represents an open operation.
type Open struct {
	Name string
}

func (o Open) Path() string { return o.Name }

// OpenFile represents an open with flags operation.
type OpenFile struct {
	Name string
	Flag int
	Perm fs.FileMode
}

func (o OpenFile) Path() string { return o.Name }

// Remove represents a remove operation.
type Remove struct {
	Name string
}

func (o Remove) Path() string { return o.Name }

// RemoveAll represents a remove -r operation.
type RemoveAll struct {
	PathName string
}

func (o RemoveAll) Path() string { return o.PathName }

// Rename represents a rename/move operation.
type Rename struct {
	Oldname string
	Newname string
}

func (o Rename) Path() string { return o.Oldname }

// Stat represents a stat operation.
type Stat struct {
	Name string
}

func (o Stat) Path() string { return o.Name }

// Readdir represents a readdir operation on a directory.
type Readdir struct {
	Name  string
	Count int
}

func (o Readdir) Path() string { return o.Name }

// Readdirnames represents a readdirnames operation on a directory.
type Readdirnames struct {
	Name  string
	Count int
}

func (o Readdirnames) Path() string { return o.Name }
