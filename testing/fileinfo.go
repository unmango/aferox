package testing

import (
	"io/fs"
	"time"
)

type FileInfo struct {
	IsDirValue   bool
	ModTimeValue time.Time
	ModeValue    fs.FileMode
	NameValue    string
	SizeValue    int64
	SysValue     any
}

// IsDir implements fs.FileInfo.
func (f *FileInfo) IsDir() bool {
	return f.IsDirValue
}

// ModTime implements fs.FileInfo.
func (f *FileInfo) ModTime() time.Time {
	return f.ModTimeValue
}

// Mode implements fs.FileInfo.
func (f *FileInfo) Mode() fs.FileMode {
	return f.ModeValue
}

// Name implements fs.FileInfo.
func (f *FileInfo) Name() string {
	return f.NameValue
}

// Size implements fs.FileInfo.
func (f *FileInfo) Size() int64 {
	return f.SizeValue
}

// Sys implements fs.FileInfo.
func (f *FileInfo) Sys() any {
	return f.SysValue
}

var _ fs.FileInfo = (*FileInfo)(nil)
