package sync

import (
	"context"

	"github.com/spf13/afero"
	"github.com/unmango/go/fs/internal"
	"github.com/unmango/go/option"
)

type (
	Func     func(context.Context) error
	Strategy func(context.Context, afero.Fs, afero.Fs) error
)

var DefaultStrategy Strategy = CopyStrategy

type Option func(*Fs)

type Fs struct {
	afero.Fs
	buf      afero.Fs
	strategy Strategy
}

func (f *Fs) Sync(ctx context.Context) error {
	return f.strategy(ctx, f.buf, f.Fs)
}

func NewFs(base, buf afero.Fs, options ...Option) (afero.Fs, Func) {
	fs := &Fs{base, buf, DefaultStrategy}
	option.ApplyAll(fs, options)
	return fs, fs.Sync
}

func CopyStrategy(_ context.Context, src, dest afero.Fs) error {
	return internal.Copy(src, dest)
}

func WithStrategy(strategy Strategy) Option {
	return func(f *Fs) {
		f.strategy = strategy
	}
}
