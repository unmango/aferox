package docker

import (
	"io/fs"
	"time"

	"github.com/docker/docker/api/types/container"
)

type FileInfo struct {
	stat container.PathStat
}

// IsDir implements fs.FileInfo.
func (f FileInfo) IsDir() bool {
	return f.stat.Mode.IsDir()
}

// ModTime implements fs.FileInfo.
func (f FileInfo) ModTime() time.Time {
	return f.stat.Mtime
}

// Mode implements fs.FileInfo.
func (f FileInfo) Mode() fs.FileMode {
	return f.stat.Mode
}

// Name implements fs.FileInfo.
func (f FileInfo) Name() string {
	return f.stat.Name
}

// Size implements fs.FileInfo.
func (f FileInfo) Size() int64 {
	return f.stat.Size
}

// Sys implements fs.FileInfo.
func (f FileInfo) Sys() any {
	return f.stat
}
