package content

import (
	"bytes"
	"io/fs"
	"syscall"

	"github.com/google/go-github/v82/github"
	"github.com/unmango/aferox/github/ghpath"
	"github.com/unmango/aferox/github/internal"
)

type File struct {
	internal.ReadOnlyFile
	ghpath.ContentPath

	client  *github.Client
	content *github.RepositoryContent

	reader *bytes.Buffer
}

// Close implements afero.File.
func (f *File) Close() error {
	return nil
}

// Name implements afero.File.
func (f *File) Name() string {
	return f.content.GetPath()
}

// Read implements afero.File.
func (f *File) Read(p []byte) (n int, err error) {
	if err = f.ensure(); err != nil {
		return
	}

	return f.reader.Read(p)
}

// ReadAt implements afero.File.
func (f *File) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, syscall.EROFS
}

// Readdir implements afero.File.
func (f *File) Readdir(count int) ([]fs.FileInfo, error) {
	return nil, syscall.ENOTDIR
}

// Readdirnames implements afero.File.
func (f *File) Readdirnames(n int) ([]string, error) {
	return nil, syscall.ENOTDIR
}

// Seek implements afero.File.
func (f *File) Seek(offset int64, whence int) (int64, error) {
	return 0, syscall.EPERM
}

// Stat implements afero.File.
func (f *File) Stat() (fs.FileInfo, error) {
	return &FileInfo{content: f.content}, nil
}

func (f *File) ensure() error {
	if f.reader != nil {
		return nil
	}

	content, err := f.content.GetContent()
	if err != nil {
		return err
	}

	f.reader = bytes.NewBufferString(content)
	return nil
}
