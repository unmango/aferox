package github

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/google/go-github/v72/github"
	"github.com/spf13/afero"
	"github.com/unmango/aferox/github/ghpath"
	"github.com/unmango/aferox/github/internal"
	"github.com/unmango/aferox/github/user"
)

type Client = github.Client

var NewClient = github.NewClient

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
		return user.Open(context.TODO(), f.client, path)
	}
}

// OpenFile implements afero.Fs.
func (f *Fs) OpenFile(name string, _ int, _ fs.FileMode) (afero.File, error) {
	if path, err := ghpath.Parse(name); err != nil {
		return nil, fmt.Errorf("open %s: %w", name, err)
	} else {
		return user.Open(context.TODO(), f.client, path)
	}
}

// Stat implements afero.Fs.
func (f *Fs) Stat(name string) (fs.FileInfo, error) {
	if path, err := ghpath.Parse(name); err != nil {
		return nil, fmt.Errorf("stat %s: %w", name, err)
	} else {
		return user.Stat(context.TODO(), f.client, path)
	}
}

func NewFs(gh *github.Client) afero.Fs {
	if gh == nil {
		gh = internal.DefaultClient()
	}

	return &Fs{client: gh}
}
