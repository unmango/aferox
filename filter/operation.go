package filter

import (
	"io/fs"
	"time"
)

// Operation represents a filesystem operation being performed.
// Use a type switch to determine the specific operation type.
type Operation interface {
	operation()
	Path() string
}

// ChmodOp represents a chmod operation.
type ChmodOp struct {
	Name string
	Mode fs.FileMode
}

func (ChmodOp) operation()     {}
func (o ChmodOp) Path() string { return o.Name }

// ChownOp represents a chown operation.
type ChownOp struct {
	Name string
	UID  int
	GID  int
}

func (ChownOp) operation()     {}
func (o ChownOp) Path() string { return o.Name }

// ChtimesOp represents a chtimes operation.
type ChtimesOp struct {
	Name  string
	Atime time.Time
	Mtime time.Time
}

func (ChtimesOp) operation()     {}
func (o ChtimesOp) Path() string { return o.Name }

// CreateOp represents a create operation.
type CreateOp struct {
	Name string
}

func (CreateOp) operation()     {}
func (o CreateOp) Path() string { return o.Name }

// MkdirOp represents a mkdir operation.
type MkdirOp struct {
	Name string
	Perm fs.FileMode
}

func (MkdirOp) operation()     {}
func (o MkdirOp) Path() string { return o.Name }

// MkdirAllOp represents a mkdir -p operation.
type MkdirAllOp struct {
	PathName string
	Perm     fs.FileMode
}

func (MkdirAllOp) operation()     {}
func (o MkdirAllOp) Path() string { return o.PathName }

// OpenOp represents an open operation.
type OpenOp struct {
	Name string
}

func (OpenOp) operation()     {}
func (o OpenOp) Path() string { return o.Name }

// OpenFileOp represents an open with flags operation.
type OpenFileOp struct {
	Name string
	Flag int
	Perm fs.FileMode
}

func (OpenFileOp) operation()     {}
func (o OpenFileOp) Path() string { return o.Name }

// RemoveOp represents a remove operation.
type RemoveOp struct {
	Name string
}

func (RemoveOp) operation()     {}
func (o RemoveOp) Path() string { return o.Name }

// RemoveAllOp represents a remove -r operation.
type RemoveAllOp struct {
	PathName string
}

func (RemoveAllOp) operation()     {}
func (o RemoveAllOp) Path() string { return o.PathName }

// RenameOp represents a rename/move operation.
type RenameOp struct {
	Oldname string
	Newname string
}

func (RenameOp) operation()     {}
func (o RenameOp) Path() string { return o.Oldname }

// StatOp represents a stat operation.
type StatOp struct {
	Name string
}

func (StatOp) operation()     {}
func (o StatOp) Path() string { return o.Name }

// ReaddirOp represents a readdir operation on a directory.
type ReaddirOp struct {
	Name  string
	Count int
}

func (ReaddirOp) operation()     {}
func (o ReaddirOp) Path() string { return o.Name }

// ReaddirnamesOp represents a readdirnames operation on a directory.
type ReaddirnamesOp struct {
	Name  string
	Count int
}

func (ReaddirnamesOp) operation()     {}
func (o ReaddirnamesOp) Path() string { return o.Name }
