package docker

import (
	"context"

	ctr "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func Stat(
	ctx context.Context,
	client client.ContainerAPIClient,
	container, path string,
) (info FileInfo, err error) {
	var stat ctr.PathStat
	stat, err = client.ContainerStatPath(ctx, container, path)
	if err != nil {
		return
	}

	return FileInfo{stat}, nil
}
