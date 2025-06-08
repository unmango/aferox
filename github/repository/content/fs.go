package content

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/google/go-github/v72/github"
	"github.com/spf13/afero"
	"github.com/unmango/aferox/github/ghpath"
	"github.com/unmango/aferox/github/internal"
)

type Fs struct {
	internal.ReadOnlyFs
	ghpath.BranchPath
	client *github.Client
}

// Name implements afero.Fs.
func (f *Fs) Name() string {
	return fmt.Sprint(f.BranchPath)
}

// Open implements afero.Fs.
func (f *Fs) Open(name string) (afero.File, error) {
	if path, err := f.Parse(name); err != nil {
		return nil, fmt.Errorf("open: %w", err)
	} else {
		return Open(context.TODO(), f.client, path)
	}
}

// OpenFile implements afero.Fs.
func (f *Fs) OpenFile(name string, _ int, _ fs.FileMode) (afero.File, error) {
	if path, err := f.Parse(name); err != nil {
		return nil, fmt.Errorf("open: %w", err)
	} else {
		return Open(context.TODO(), f.client, path)
	}
}

// Stat implements afero.Fs.
func (f *Fs) Stat(name string) (fs.FileInfo, error) {
	if path, err := f.Parse(name); err != nil {
		return nil, fmt.Errorf("stat: %w", err)
	} else {
		return Stat(context.TODO(), f.client, path)
	}
}

func Open(ctx context.Context, client *github.Client, path ghpath.Path) (afero.File, error) {
	content, err := ghpath.ParseContent(path)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	file, dir, _, err := client.Repositories.GetContents(ctx,
		content.Owner,
		content.Repository,
		content.Content,
		&github.RepositoryContentGetOptions{
			Ref: content.Branch,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	if file != nil {
		return &File{
			ContentPath: content,
			client:      client,
			content:     file,
		}, nil
	} else {
		return &Directory{
			ContentPath: content,
			client:      client,
			content:     dir,
		}, nil
	}
}

func Stat(ctx context.Context, client *github.Client, path ghpath.Path) (fs.FileInfo, error) {
	content, err := ghpath.ParseContent(path)
	if err != nil {
		return nil, fmt.Errorf("stat: %w", err)
	}

	file, dir, _, err := client.Repositories.GetContents(ctx,
		content.Owner,
		content.Repository,
		content.Content,
		&github.RepositoryContentGetOptions{
			Ref: content.Branch,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("stat: %w", err)
	}

	if file != nil {
		return &FileInfo{content: file}, nil
	} else {
		return &DirectoryInfo{
			name:    content.Content,
			content: dir,
		}, nil
	}
}

func NewFs(gh *github.Client, owner, repo, branch string) afero.Fs {
	return &Fs{
		client:     gh,
		BranchPath: ghpath.NewBranchPath(owner, repo, branch),
	}
}
