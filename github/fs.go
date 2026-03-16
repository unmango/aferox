package github

import (
	"io/fs"

	gogithub "github.com/google/go-github/v84/github"
	"github.com/spf13/afero"
	"github.com/unmango/aferox/github/internal"
	"github.com/unstoppablemango/ihfs/ghfs"
)

type Client = gogithub.Client

var NewClient = gogithub.NewClient

type Fs struct {
	internal.ReadOnlyFs
	ghfs *ghfs.Fs
}

// Name implements afero.Fs.
func (f *Fs) Name() string {
	return f.ghfs.Name()
}

// Open implements afero.Fs.
func (f *Fs) Open(name string) (afero.File, error) {
	file, err := f.ghfs.Open(name)
	if err != nil {
		return nil, err
	}
	return &File{name: name, file: file}, nil
}

// OpenFile implements afero.Fs.
func (f *Fs) OpenFile(name string, _ int, _ fs.FileMode) (afero.File, error) {
	return f.Open(name)
}

// Stat implements afero.Fs.
func (f *Fs) Stat(name string) (fs.FileInfo, error) {
	file, err := f.ghfs.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return file.Stat()
}

func NewFs(gh *gogithub.Client) afero.Fs {
	if gh == nil {
		gh = internal.DefaultClient()
	}
	return &Fs{ghfs: ghfs.New(ghfs.WithClient(gh))}
}
