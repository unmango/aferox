package aferox

import (
	"fmt"

	"github.com/spf13/afero"
	"github.com/unmango/aferox/context"
)

type contextKey struct{}

var defaultContextFs = afero.NewOsFs()

func SetContext(fs afero.Fs, ctx context.Context) error {
	if ctxfs, ok := fs.(context.Setter); !ok {
		return fmt.Errorf("context not supported: %s", fs.Name())
	} else {
		ctxfs.SetContext(ctx)
	}

	return nil
}

func FromContext(ctx context.Context) afero.Fs {
	if fs := ctx.Value(contextKey{}); fs != nil {
		return fs.(afero.Fs)
	} else {
		return defaultContextFs
	}
}
