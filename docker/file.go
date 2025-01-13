package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/unmango/aferox/docker/internal"
)

type File struct {
	client    client.ContainerAPIClient
	container string
	name      string

	close  func() error
	stat   container.PathStat
	reader io.Reader
}

// Close implements afero.File.
func (f *File) Close() error {
	if f.close != nil {
		return f.close()
	}

	return nil
}

// Name implements afero.File.
func (f *File) Name() string {
	return f.name
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
	if ra, ok := f.reader.(io.ReaderAt); ok {
		return ra.ReadAt(p, off)
	}

	return 0, syscall.EPERM
}

// Readdir implements afero.File.
func (f *File) Readdir(count int) ([]fs.FileInfo, error) {
	ctx := context.TODO()
	buf := &bytes.Buffer{}
	err := f.execo(ctx, internal.ExecOptions{
		Cmd:    []string{"dir", "-x1", f.name},
		Stdout: buf,
	})
	if err != nil {
		return nil, err
	}

	cleaned := strings.TrimSpace(buf.String())
	paths := strings.Split(cleaned, "\n")
	length := min(len(paths), count)
	infos := make([]fs.FileInfo, length)

	for i := 0; i < length; i++ {
		stat, err := Stat(ctx, f.client, f.container, paths[i])
		if err != nil {
			return nil, err
		}

		infos[i] = stat
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
	panic("unimplemented")
}

// Stat implements afero.File.
func (f *File) Stat() (fs.FileInfo, error) {
	return Stat(context.TODO(), f.client, f.container, f.name)
}

// Sync implements afero.File.
func (f *File) Sync() error {
	// TODO: Memory sync vs container fs sync?
	return f.exec(context.TODO(), "sync", f.name)
}

// Truncate implements afero.File.
func (f *File) Truncate(size int64) error {
	return f.exec(context.TODO(),
		"truncate", fmt.Sprintf("--size=%d", size), f.name,
	)
}

// Write implements afero.File.
func (f *File) Write(p []byte) (n int, err error) {
	// TODO: This only works in one go
	content := &bytes.Buffer{}
	w := tar.NewWriter(content)
	err = w.WriteHeader(&tar.Header{
		Name: filepath.Base(f.name),
		Size: int64(len(p)),
	})
	if err != nil {
		return
	}

	err = f.client.CopyToContainer(context.TODO(),
		f.container,
		f.name,
		bytes.NewBuffer(p),
		container.CopyToContainerOptions{},
	)
	if err != nil {
		return
	}

	n = len(p)
	return
}

// WriteAt implements afero.File.
func (f *File) WriteAt(p []byte, off int64) (n int, err error) {
	panic("unimplemented")
}

// WriteString implements afero.File.
func (f *File) WriteString(s string) (ret int, err error) {
	content := &bytes.Buffer{}
	w := tar.NewWriter(content)
	err = w.WriteHeader(&tar.Header{
		Name: filepath.Base(f.name),
		Size: int64(len(s)),
	})
	if err != nil {
		return
	}

	ret, err = io.WriteString(w, s)
	if err != nil {
		return
	}

	err = f.client.CopyToContainer(context.TODO(),
		f.container,
		filepath.Dir(f.name),
		content,
		container.CopyToContainerOptions{},
	)
	if err != nil {
		return
	}

	return
}

func (f *File) ensure() error {
	if f.reader != nil {
		return nil
	}

	ctx := context.TODO()
	reader, stat, err := f.client.CopyFromContainer(ctx,
		f.container,
		f.name,
	)
	if err != nil {
		return err
	}

	tar := tar.NewReader(reader)
	if _, err = tar.Next(); err != nil {
		return err
	}

	f.reader = tar
	f.stat = stat
	f.close = reader.Close

	return nil
}

func (f *File) exec(ctx context.Context, cmd ...string) error {
	return f.execo(ctx, internal.ExecOptions{
		Cmd: cmd,
	})
}

func (f *File) execo(ctx context.Context, options internal.ExecOptions) error {
	return internal.Exec(ctx, f.client, f.container, options)
}
