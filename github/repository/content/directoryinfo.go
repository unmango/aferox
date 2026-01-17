package content

import (
	"io/fs"
	"os"
	"time"

	"github.com/google/go-github/v81/github"
)

type DirectoryInfo struct {
	name    string
	content []*github.RepositoryContent
}

// IsDir implements fs.FileInfo.
func (d *DirectoryInfo) IsDir() bool {
	return true
}

// ModTime implements fs.FileInfo.
func (d *DirectoryInfo) ModTime() time.Time {
	panic("unimplemented")
}

// Mode implements fs.FileInfo.
func (d *DirectoryInfo) Mode() fs.FileMode {
	return os.ModeDir
}

// Name implements fs.FileInfo.
func (d *DirectoryInfo) Name() string {
	return d.name
}

// Size implements fs.FileInfo.
func (d *DirectoryInfo) Size() (s int64) {
	for _, c := range d.content {
		s += int64(c.GetSize())
	}

	return
}

// Sys implements fs.FileInfo.
func (d *DirectoryInfo) Sys() any {
	return d.content
}
