package release

import (
	"io/fs"
	"os"
	"time"

	"github.com/google/go-github/v80/github"
)

type FileInfo struct {
	release *github.RepositoryRelease
}

// IsDir implements fs.FileInfo.
func (f *FileInfo) IsDir() bool {
	return true
}

// ModTime implements fs.FileInfo.
func (f *FileInfo) ModTime() time.Time {
	return f.release.GetCreatedAt().Time
}

// Mode implements fs.FileInfo.
func (f *FileInfo) Mode() fs.FileMode {
	return os.ModeDir
}

// Name implements fs.FileInfo.
func (f *FileInfo) Name() string {
	return f.release.GetName()
}

// Size implements fs.FileInfo.
func (f *FileInfo) Size() int64 {
	panic("unimplemented")
}

// Sys implements fs.FileInfo.
func (f *FileInfo) Sys() any {
	return f.release
}
