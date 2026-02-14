package asset

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"syscall"

	"github.com/google/go-github/v82/github"
	"github.com/unmango/aferox/github/ghpath"
	"github.com/unmango/aferox/github/internal"
)

type File struct {
	internal.ReadOnlyFile
	ghpath.ReleasePath

	client *github.Client
	asset  *github.ReleaseAsset

	reader io.ReadCloser
}

// Close implements afero.File.
func (f *File) Close() error {
	if f.reader != nil {
		return f.reader.Close()
	} else {
		return nil
	}
}

// Name implements afero.File.
func (f *File) Name() string {
	return f.asset.GetName()
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
	return 0, syscall.EPERM
}

// Readdir implements afero.File.
func (f *File) Readdir(count int) (infos []fs.FileInfo, err error) {
	if !f.isArchive() {
		return nil, syscall.ENOTDIR
	}

	if err := f.ensure(); err != nil {
		return nil, err
	}

	r := f.reader
	if f.isGzip() {
		if r, err = gzip.NewReader(r); err != nil {
			return
		}
	}

	tar := tar.NewReader(r)
	for {
		h, err := tar.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading archive: %w", err)
		}

		// TODO: Handle directories
		infos = append(infos, h.FileInfo())
	}

	return infos, nil
}

// Readdirnames implements afero.File.
func (f *File) Readdirnames(n int) ([]string, error) {
	infos, err := f.Readdir(n)
	if err != nil {
		return nil, err
	}

	length := min(n, len(infos))
	names := make([]string, length)
	for i, info := range infos {
		names[i] = info.Name()
	}

	return names, nil
}

// Seek implements afero.File.
func (f *File) Seek(offset int64, whence int) (int64, error) {
	return 0, syscall.EPERM
}

// Stat implements afero.File.
func (f *File) Stat() (fs.FileInfo, error) {
	return &FileInfo{asset: f.asset}, nil
}

func (f *File) ensure() error {
	if f.reader != nil {
		return nil
	}

	reader, _, err := f.client.Repositories.DownloadReleaseAsset(
		context.TODO(),
		f.Owner,
		f.Repository,
		f.asset.GetID(),
		http.DefaultClient,
	)
	if err != nil {
		return err
	}

	f.reader = reader
	return nil
}

func (f *File) isArchive() bool {
	name := f.asset.GetName()
	return strings.HasSuffix(name, ".tar.gz") ||
		strings.HasSuffix(name, ".tar")
}

func (f *File) isGzip() bool {
	return strings.HasSuffix(f.asset.GetName(), ".gz")
}
