package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"syscall"

	"github.com/google/go-github/v68/github"
	"github.com/unmango/aferox/github/ghpath"
	"github.com/unmango/aferox/github/internal"
	"github.com/unmango/aferox/github/repository/release"
)

type File struct {
	internal.ReadOnlyFile
	ghpath.OwnerPath

	client *github.Client
	repo   *github.Repository

	reader *bytes.Buffer
}

// Close implements afero.File.
func (f *File) Close() error {
	return nil
}

// Name implements afero.File.
func (f *File) Name() string {
	return f.repo.GetName()
}

// Read implements afero.File.
func (f *File) Read(p []byte) (n int, err error) {
	if err = f.ensure(); err != nil {
		return
	} else {
		return f.reader.Read(p)
	}
}

// ReadAt implements afero.File.
func (f *File) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, syscall.EPERM
}

// Readdir implements afero.File.
func (f *File) Readdir(count int) ([]fs.FileInfo, error) {
	return release.Readdir(context.TODO(),
		f.client,
		f.Owner,
		f.repo.GetName(),
		count,
	)
}

// Readdirnames implements afero.File.
func (f *File) Readdirnames(n int) ([]string, error) {
	return release.Readdirnames(context.TODO(),
		f.client,
		f.Owner,
		f.repo.GetName(),
		n,
	)
}

// Seek implements afero.File.
func (f *File) Seek(offset int64, whence int) (int64, error) {
	return 0, syscall.EPERM
}

// Stat implements afero.File.
func (f *File) Stat() (fs.FileInfo, error) {
	return &FileInfo{repo: f.repo}, nil
}

func (f *File) ensure() error {
	if f.reader != nil {
		return nil
	}

	data, err := json.Marshal(f.repo)
	if err != nil {
		return fmt.Errorf("marshaling repo: %w", err)
	}

	f.reader = bytes.NewBuffer(data)
	return nil
}
