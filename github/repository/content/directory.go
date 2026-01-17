package content

import (
	"io/fs"
	"syscall"

	"github.com/google/go-github/v81/github"
	"github.com/unmango/aferox/github/ghpath"
	"github.com/unmango/aferox/github/internal"
)

type Directory struct {
	internal.ReadOnlyFile
	ghpath.ContentPath

	client  *github.Client
	content []*github.RepositoryContent
}

// Close implements afero.File.
func (d *Directory) Close() error {
	return nil
}

// Name implements afero.File.
func (d *Directory) Name() string {
	return d.Content
}

// Read implements afero.File.
func (d *Directory) Read(p []byte) (n int, err error) {
	return 0, syscall.EISDIR
}

// ReadAt implements afero.File.
func (d *Directory) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, syscall.EISDIR
}

// Readdir implements afero.File.
func (d *Directory) Readdir(count int) ([]fs.FileInfo, error) {
	length := min(count, len(d.content))
	files := make([]fs.FileInfo, length)

	for i := 0; i < length; i++ {
		files[i] = &FileInfo{content: d.content[i]}
	}

	return files, nil
}

// Readdirnames implements afero.File.
func (d *Directory) Readdirnames(n int) ([]string, error) {
	infos, err := d.Readdir(n)
	if err != nil {
		return nil, err
	}

	names := []string{}
	for _, info := range infos {
		names = append(names, info.Name())
	}

	return names, nil
}

// Seek implements afero.File.
func (d *Directory) Seek(offset int64, whence int) (int64, error) {
	return 0, syscall.EISDIR
}

// Stat implements afero.File.
func (d *Directory) Stat() (fs.FileInfo, error) {
	return &DirectoryInfo{
		name:    d.Name(),
		content: d.content,
	}, nil
}
