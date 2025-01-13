package docker

import (
	"fmt"
	"io/fs"
	"time"

	"github.com/docker/docker/client"
	"github.com/spf13/afero"
	"github.com/unmango/aferox/context"
	"github.com/unmango/aferox/docker/internal"
)

type Fs struct {
	client    client.ContainerAPIClient
	container string
}

// Chmod implements afero.Fs.
func (f Fs) Chmod(ctx context.Context, name string, mode fs.FileMode) error {
	return f.exec(ctx, "chmod", mode.String(), name)
}

// Chown implements afero.Fs.
func (f Fs) Chown(ctx context.Context, name string, uid int, gid int) error {
	return f.exec(ctx,
		"chown", fmt.Sprintf("%d:%d", uid, gid), name,
	)
}

// Chtimes implements afero.Fs.
func (f Fs) Chtimes(ctx context.Context, name string, atime time.Time, mtime time.Time) error {
	panic("unimplemented")
}

// Create implements afero.Fs.
func (f Fs) Create(ctx context.Context, name string) (afero.File, error) {
	err := f.exec(ctx, "touch", name)
	if err != nil {
		return nil, err
	}

	// TODO: Less lazy?
	return &File{
		client:    f.client,
		container: f.container,
		name:      name,
	}, nil
}

// Mkdir implements afero.Fs.
func (f Fs) Mkdir(ctx context.Context, name string, perm fs.FileMode) error {
	return f.exec(ctx,
		"mkdir", fmt.Sprintf("--mode=%d", perm), name,
	)
}

// MkdirAll implements afero.Fs.
func (f Fs) MkdirAll(ctx context.Context, path string, perm fs.FileMode) error {
	return f.exec(ctx,
		"mkdir", "--parents", fmt.Sprintf("--mode=%d", perm), path,
	)
}

// Name implements afero.Fs.
func (f Fs) Name() string {
	return f.container
}

// Open implements afero.Fs.
func (f Fs) Open(ctx context.Context, name string) (afero.File, error) {
	// TODO: Less lazy?
	return &File{
		client:    f.client,
		container: f.container,
		name:      name,
	}, nil
}

// OpenFile implements afero.Fs.
func (f Fs) OpenFile(ctx context.Context, name string, flag int, perm fs.FileMode) (afero.File, error) {
	// TODO: Actual implementation
	// TODO: Less lazy?
	return &File{
		client:    f.client,
		container: f.container,
		name:      name,
	}, nil
}

// Remove implements afero.Fs.
func (f Fs) Remove(ctx context.Context, name string) error {
	return f.exec(ctx, "rm", name)
}

// RemoveAll implements afero.Fs.
func (f Fs) RemoveAll(ctx context.Context, path string) error {
	return f.exec(ctx, "rm", "--recursive", path)
}

// Rename implements afero.Fs.
func (f Fs) Rename(ctx context.Context, oldname string, newname string) error {
	return f.exec(ctx, "mv", oldname, newname)
}

// Stat implements afero.Fs.
func (f Fs) Stat(ctx context.Context, name string) (fs.FileInfo, error) {
	return Stat(ctx, f.client, f.container, name)
}

func (f Fs) exec(ctx context.Context, cmd ...string) error {
	return internal.Exec(ctx, f.client, f.container,
		internal.ExecOptions{
			Cmd: cmd,
		},
	)
}

func NewFs(client client.ContainerAPIClient, container string) context.Fs {
	return Fs{client, container}
}
