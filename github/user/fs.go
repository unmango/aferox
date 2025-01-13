package user

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/google/go-github/v68/github"
	"github.com/spf13/afero"
	"github.com/unmango/aferox/github/ghpath"
	"github.com/unmango/aferox/github/internal"
	"github.com/unmango/aferox/github/repository"
)

type Fs struct {
	internal.ReadOnlyFs
	client *github.Client
}

// Name implements afero.Fs.
func (g *Fs) Name() string {
	return "https://github.com"
}

// Open implements afero.Fs.
func (f *Fs) Open(name string) (afero.File, error) {
	if path, err := ghpath.Parse(name); err != nil {
		return nil, fmt.Errorf("open %s: %w", name, err)
	} else {
		return Open(context.TODO(), f.client, path)
	}
}

// OpenFile implements afero.Fs.
func (f *Fs) OpenFile(name string, _ int, _ fs.FileMode) (afero.File, error) {
	if path, err := ghpath.Parse(name); err != nil {
		return nil, fmt.Errorf("open %s: %w", name, err)
	} else {
		return Open(context.TODO(), f.client, path)
	}
}

// Stat implements afero.Fs.
func (f *Fs) Stat(name string) (fs.FileInfo, error) {
	if path, err := ghpath.Parse(name); err != nil {
		return nil, fmt.Errorf("stat %s: %w", name, err)
	} else {
		return Stat(context.TODO(), f.client, path)
	}
}

func NewFs(gh *github.Client) afero.Fs {
	return &Fs{client: gh}
}

func Open(ctx context.Context, gh *github.Client, path ghpath.Path) (afero.File, error) {
	if _, err := path.Repository(); err == nil {
		return repository.Open(ctx, gh, path)
	}

	owner, err := ghpath.ParseOwner(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}

	user, _, err := gh.Users.Get(ctx, owner.Owner)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}

	return &File{
		client: gh,
		user:   user,
	}, nil
}

func Stat(ctx context.Context, gh *github.Client, path ghpath.Path) (fs.FileInfo, error) {
	if _, err := path.Repository(); err == nil {
		return repository.Stat(ctx, gh, path)
	}

	owner, err := ghpath.ParseOwner(path)
	if err != nil {
		return nil, fmt.Errorf("stat %s: %w", path, err)
	}

	user, _, err := gh.Users.Get(ctx, owner.Owner)
	if err != nil {
		return nil, fmt.Errorf("stat %s: %w", path, err)
	}

	return &FileInfo{
		client: gh,
		user:   user,
	}, nil
}
