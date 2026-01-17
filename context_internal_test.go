package aferox

import (
	stdctx "context"
	"testing"

	"github.com/spf13/afero"
)

// Test FromContext with a value in context using the private contextKey
func TestFromContextWithValue(t *testing.T) {
	expectedFs := afero.NewMemMapFs()
	ctx := stdctx.WithValue(stdctx.Background(), contextKey{}, expectedFs)

	fs := FromContext(ctx)

	if fs != expectedFs {
		t.Errorf("Expected fs to be %v, got %v", expectedFs, fs)
	}
}
