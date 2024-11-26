package docker

import (
	"archive/tar"
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/client"
	"github.com/spf13/afero"
	"github.com/spf13/afero/tarfs"
)

type ReadOnlyFs struct {
	*tarfs.Fs
	client    client.ContainerAPIClient
	container string
}

func (f *ReadOnlyFs) Name() string {
	return fmt.Sprintf("ReadOnly: %s", f.container)
}

func (f *ReadOnlyFs) Open(name string) (afero.File, error) {
	if err := f.ensure(); err != nil {
		return nil, fmt.Errorf("open %s: %w", name, err)
	}

	return f.Fs.Open(name)
}

func (f *ReadOnlyFs) OpenFile(name string, _ int, _ os.FileMode) (afero.File, error) {
	if err := f.ensure(); err != nil {
		return nil, fmt.Errorf("open %s: %w", name, err)
	}

	return f.Fs.Open(name)
}

func (f *ReadOnlyFs) Stat(name string) (os.FileInfo, error) {
	if err := f.ensure(); err != nil {
		return nil, fmt.Errorf("stat %s: %w", name, err)
	}

	return f.Fs.Stat(name)
}

func (f *ReadOnlyFs) ensure() error {
	if f.Fs != nil {
		return nil
	}

	reader, _, err := f.client.CopyFromContainer(context.TODO(), f.container, "")
	if err != nil {
		return err
	}

	f.Fs = tarfs.New(tar.NewReader(reader))
	return nil
}

func NewReadOnly(client client.ContainerAPIClient, container string) afero.Fs {
	return &ReadOnlyFs{
		client:    client,
		container: container,
	}
}
