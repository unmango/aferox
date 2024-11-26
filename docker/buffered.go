package docker

import (
	"archive/tar"
	"bytes"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/afero"
	"github.com/unmango/go/fs/context"
	"github.com/unmango/go/fs/sync"
)

func NewBufferedFs(client client.ContainerAPIClient, ctr string) (afero.Fs, sync.Func) {
	buf := afero.NewMemMapFs()
	base := afero.NewCopyOnWriteFs(
		context.TodoFs(NewFs(client, ctr)),
		buf,
	)

	strategy := func(ctx context.Context, src, _ afero.Fs) error {
		buf := &bytes.Buffer{}
		archive := tar.NewWriter(buf)
		if err := archive.AddFS(afero.NewIOFS(src)); err != nil {
			return fmt.Errorf("adding src to archive: %w", err)
		}

		return client.CopyToContainer(ctx, ctr, "", buf,
			container.CopyToContainerOptions{},
		)
	}

	return sync.NewFs(base, buf, sync.WithStrategy(strategy))
}
