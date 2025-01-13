package aferox

import (
	"io"

	"github.com/docker/docker/client"
	"github.com/spf13/afero"
	"github.com/unmango/aferox/context"
	"github.com/unmango/aferox/docker"
	"github.com/unmango/aferox/github"
	"github.com/unmango/aferox/github/repository"
	"github.com/unmango/aferox/github/repository/content"
	"github.com/unmango/aferox/github/repository/release"
	"github.com/unmango/aferox/github/user"
	"github.com/unmango/aferox/writer"
)

type (
	Docker                  = docker.Fs
	GitHub                  = github.Fs
	GitHubRelease           = release.Fs
	GitHubRepository        = repository.Fs
	GitHubRepositoryContent = content.Fs
	GitHubUser              = user.Fs
	Writer                  = writer.Fs
)

func NewWriter(w io.Writer) afero.Fs {
	return writer.NewFs(w)
}

func NewDocker(client client.ContainerAPIClient, container string) context.Fs {
	return docker.NewFs(client, container)
}
