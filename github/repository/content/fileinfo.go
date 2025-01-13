package content

import (
	"io/fs"
	"os"
	"time"

	"github.com/google/go-github/v68/github"
)

type FileInfo struct {
	content *github.RepositoryContent
}

// IsDir implements fs.FileInfo.
func (f *FileInfo) IsDir() bool {
	return f.content == nil
}

// ModTime implements fs.FileInfo.
func (f *FileInfo) ModTime() time.Time {
	panic("unimplemented")
}

// Mode implements fs.FileInfo.
func (f *FileInfo) Mode() fs.FileMode {
	return os.ModePerm
}

// Name implements fs.FileInfo.
func (f *FileInfo) Name() string {
	return f.content.GetName()
}

// Size implements fs.FileInfo.
func (f *FileInfo) Size() int64 {
	return int64(f.content.GetSize())
}

// Sys implements fs.FileInfo.
func (f *FileInfo) Sys() any {
	return f.content
}
