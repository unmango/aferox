package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"

	"github.com/google/go-github/v75/github"
	"github.com/unmango/aferox/github/internal"
	"github.com/unmango/aferox/github/repository"
)

type File struct {
	internal.ReadOnlyFile
	client *github.Client
	user   *github.User

	reader *bytes.Reader
}

// Close implements afero.File.
func (f *File) Close() error {
	return nil
}

// Name implements afero.File.
func (f *File) Name() string {
	return f.user.GetLogin()
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
	if err = f.ensure(); err != nil {
		return
	} else {
		return f.reader.ReadAt(p, off)
	}
}

// Readdir implements afero.File.
func (f *File) Readdir(count int) ([]fs.FileInfo, error) {
	return repository.Readdir(context.TODO(), f.client, f.Name(), count)
}

// Readdirnames implements afero.File.
func (f *File) Readdirnames(n int) ([]string, error) {
	return repository.Readdirnames(context.TODO(), f.client, f.Name(), n)
}

// Seek implements afero.File.
func (f *File) Seek(offset int64, whence int) (int64, error) {
	if err := f.ensure(); err != nil {
		return 0, err
	} else {
		return f.reader.Seek(offset, whence)
	}
}

// Stat implements afero.File.
func (f *File) Stat() (fs.FileInfo, error) {
	return &FileInfo{
		client: f.client,
		user:   f.user,
	}, nil
}

func (f *File) ensure() error {
	if f.reader != nil {
		return nil
	}

	data, err := json.Marshal(f.user)
	if err != nil {
		return fmt.Errorf("marshaling user: %w", err)
	}

	f.reader = bytes.NewReader(data)
	return nil
}
