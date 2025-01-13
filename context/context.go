package context

import (
	"context"

	"github.com/spf13/afero"
)

type Setter interface {
	SetContext(context.Context)
}

type Accessor interface {
	Context() context.Context
}

type Context = context.Context

var (
	Background        = context.Background
	TODO              = context.TODO
	WithCancel        = context.WithCancel
	WithCancelCause   = context.WithCancelCause
	WithDeadline      = context.WithDeadline
	WithDeadlineCause = context.WithDeadlineCause
	WithValue         = context.WithValue
	WithTimeout       = context.WithTimeout
	WithTimeoutCause  = context.WithTimeoutCause
	WithoutCancel     = context.WithoutCancel
)

type union struct {
	afero.Fs
	AdapterFs
}

func NewFs(base Fs, accessor Accessor) afero.Fs {
	return &AccessorFs{accessor, base}
}

func BackgroundFs(base Fs) afero.Fs {
	return NewFs(base, AccessorFunc(Background))
}

func TodoFs(base Fs) afero.Fs {
	return NewFs(base, AccessorFunc(TODO))
}

func Discard(fs afero.Fs) AferoFs {
	return &DiscardFs{fs}
}

func Adapt(fs Fs, accessor Accessor) AferoFs {
	return &union{
		Fs:        NewFs(fs, accessor),
		AdapterFs: AdapterFs{fs},
	}
}
