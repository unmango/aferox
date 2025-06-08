package repository

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/google/go-github/v72/github"
	"github.com/spf13/afero"
	"github.com/unmango/aferox/github/ghpath"
	"github.com/unmango/aferox/github/internal"
	"github.com/unmango/aferox/github/repository/content"
	"github.com/unmango/aferox/github/repository/release"
)

type Fs struct {
	internal.ReadOnlyFs
	ghpath.OwnerPath
	client *github.Client
}

// Name implements afero.Fs.
func (f *Fs) Name() string {
	return fmt.Sprint(f.OwnerPath)
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

func NewFs(gh *github.Client, owner string) afero.Fs {
	return &Fs{
		client:    gh,
		OwnerPath: ghpath.NewOwnerPath(owner),
	}
}

func Open(ctx context.Context, gh *github.Client, path ghpath.Path) (afero.File, error) {
	if _, err := path.Release(); err == nil {
		return release.Open(ctx, gh, path)
	}
	if _, err := path.Branch(); err == nil {
		return content.Open(ctx, gh, path)
	}

	repo, err := ghpath.ParseRepository(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path %s: %w", path, err)
	}

	r, _, err := gh.Repositories.Get(ctx, repo.Owner, repo.Repository)
	if err != nil {
		return nil, err
	}

	return &File{
		client:    gh,
		repo:      r,
		OwnerPath: repo.OwnerPath,
	}, nil
}

func Readdir(ctx context.Context, gh *github.Client, user string, count int) ([]fs.FileInfo, error) {
	// TODO: count == 0
	opt := &github.RepositoryListByUserOptions{
		ListOptions: github.ListOptions{PerPage: count},
	}

	repos, _, err := gh.Repositories.ListByUser(ctx, user, opt)
	if err != nil {
		return nil, fmt.Errorf("user %s readdir: %w", user, err)
	}

	length := min(count, len(repos))
	infos := make([]fs.FileInfo, length)

	for i := 0; i < length; i++ {
		infos[i] = &FileInfo{repo: repos[i]}
	}

	return infos, nil
}

func Readdirnames(ctx context.Context, gh *github.Client, user string, n int) ([]string, error) {
	infos, err := Readdir(ctx, gh, user, n)
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
	if _, err := path.Release(); err == nil {
		return release.Stat(ctx, gh, path)
	}
	if _, err := path.Branch(); err == nil {
		return content.Stat(ctx, gh, path)
	}

	repo, err := ghpath.ParseRepository(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path %s: %w", path, err)
	}

	r, _, err := gh.Repositories.Get(ctx, repo.Owner, repo.Repository)
	if err != nil {
		return nil, fmt.Errorf("stat: %w", err)
	}

	return &FileInfo{repo: r}, nil
}
