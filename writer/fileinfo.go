package writer

import (
	"io"
	"io/fs"
	"os"
	"time"
)

type FileInfo struct {
	io.Writer
	name string
}

// IsDir implements fs.FileInfo.
func (f *FileInfo) IsDir() bool {
	return false
}

// ModTime implements fs.FileInfo.
func (f *FileInfo) ModTime() time.Time {
	return time.Time{}
}

// Mode implements fs.FileInfo.
func (f *FileInfo) Mode() fs.FileMode {
	return os.ModePerm
}

// Name implements fs.FileInfo.
func (f *FileInfo) Name() string {
	return f.name
}

// Size implements fs.FileInfo.
func (f *FileInfo) Size() int64 {
	return -1
}

// Sys implements fs.FileInfo.
func (f *FileInfo) Sys() any {
	return f.Writer
}
