package release

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"

	"github.com/google/go-github/v72/github"
	"github.com/unmango/aferox/github/ghpath"
	"github.com/unmango/aferox/github/internal"
	"github.com/unmango/aferox/github/repository/release/asset"
)

type File struct {
	internal.ReadOnlyFile
	ghpath.RepositoryPath

	client  *github.Client
	release *github.RepositoryRelease

	reader *bytes.Reader
}

// Close implements afero.File.
func (f *File) Close() error {
	f.reader = nil
	return nil
}

// Name implements afero.File.
func (f *File) Name() string {
	return f.release.GetName()
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
	return asset.Readdir(context.TODO(),
		f.client,
		f.RepositoryPath,
		f.release.GetID(),
		count,
	)
}

// Readdirnames implements afero.File.
func (f *File) Readdirnames(n int) ([]string, error) {
	return asset.Readdirnames(context.TODO(),
		f.client,
		f.RepositoryPath,
		f.release.GetID(),
		n,
	)
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
	return &FileInfo{release: f.release}, nil
}

func (f *File) ensure() error {
	if f.reader != nil {
		return nil
	}

	data, err := json.Marshal(f.release)
	if err != nil {
		return fmt.Errorf("marshaling release: %w", err)
	}

	f.reader = bytes.NewReader(data)
	return nil
}
