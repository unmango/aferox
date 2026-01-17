package release

import (
	"context"
	"fmt"
	"io/fs"
	"os"

	"github.com/google/go-github/v81/github"
	"github.com/spf13/afero"
	"github.com/unmango/aferox/github/ghpath"
	"github.com/unmango/aferox/github/internal"
	"github.com/unmango/aferox/github/repository/release/asset"
)

type Fs struct {
	internal.ReadOnlyFs
	ghpath.RepositoryPath
	client *github.Client
}

// Name implements afero.Fs.
func (f *Fs) Name() string {
	return fmt.Sprintf("%s/releases", f.RepositoryPath)
}

// Open implements afero.Fs.
func (f *Fs) Open(name string) (afero.File, error) {
	path, err := f.Parse(name)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", name, err)
	}

	return Open(context.TODO(), f.client, path)
}

// OpenFile implements afero.Fs.
func (f *Fs) OpenFile(name string, _ int, _ fs.FileMode) (afero.File, error) {
	path, err := f.Parse(name)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", name, err)
	}

	return Open(context.TODO(), f.client, path)
}

// Stat implements afero.Fs.
func (f *Fs) Stat(name string) (fs.FileInfo, error) {
	path, err := f.Parse(name)
	if err != nil {
		return nil, fmt.Errorf("stat %s: %w", name, err)
	}

	return Stat(context.TODO(), f.client, path)
}

func NewFs(gh *github.Client, owner, repository string) afero.Fs {
	return &Fs{
		client:         gh,
		RepositoryPath: ghpath.NewRepositoryPath(owner, repository),
	}
}

func Open(ctx context.Context, gh *github.Client, path ghpath.Path) (afero.File, error) {
	release, err := ghpath.ParseRelease(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path %s: %w", path, err)
	}

	if _, err := path.Asset(); err == nil {
		return asset.Open(ctx, gh, path)
	}

	id, err := releaseId(ctx, gh, release)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}

	r, _, err := gh.Repositories.GetRelease(ctx, release.Owner, release.Repository, id)
	if err != nil {
		return nil, fmt.Errorf("open %d: %w", id, err)
	}

	return &File{
		client:         gh,
		release:        r,
		RepositoryPath: release.RepositoryPath,
	}, nil
}

func Readdir(ctx context.Context, gh *github.Client, owner, repository string, count int) ([]fs.FileInfo, error) {
	// TODO: count == 0
	opt := &github.ListOptions{PerPage: count}
	releases, _, err := gh.Repositories.ListReleases(ctx, owner, repository, opt)
	if err != nil {
		return nil, fmt.Errorf("%s/%s readdir: %w", owner, repository, err)
	}

	length := min(count, len(releases))
	results := make([]fs.FileInfo, length)

	for i := 0; i < length; i++ {
		results[i] = &FileInfo{release: releases[i]}
	}

	return results, nil
}

func Readdirnames(ctx context.Context, gh *github.Client, owner, repository string, n int) ([]string, error) {
	infos, err := Readdir(ctx, gh, owner, repository, n)
	if err != nil {
		return nil, err
	}

	names := []string{}
	for _, i := range infos {
		names = append(names, i.Name())
	}

	return names, nil
}

func Stat(ctx context.Context, gh *github.Client, path ghpath.Path) (fs.FileInfo, error) {
	if _, err := path.Asset(); err == nil {
		return asset.Stat(ctx, gh, path)
	}

	release, err := ghpath.ParseRelease(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path %s: %w", release, err)
	}

	id, err := releaseId(ctx, gh, release)
	if err != nil {
		return nil, fmt.Errorf("reading release id: %w", err)
	}

	r, _, err := gh.Repositories.GetRelease(ctx, release.Owner, release.Repository, id)
	if err != nil {
		return nil, fmt.Errorf("open %d: %w", id, err)
	}

	return &FileInfo{release: r}, nil
}

func releaseId(ctx context.Context, gh *github.Client, path ghpath.ReleasePath) (int64, error) {
	if id, ok := internal.TryGetId(path.Release); ok {
		return id, nil
	}

	releases, _, err := gh.Repositories.ListReleases(ctx, path.Owner, path.Repository, nil)
	if err != nil {
		return 0, err
	}

	for _, r := range releases {
		if r.GetName() == path.Release {
			return r.GetID(), nil
		}
	}

	return 0, os.ErrNotExist
}
