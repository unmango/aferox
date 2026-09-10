package internal

import (
	"syscall"

	"github.com/spf13/afero"
)

type ReadOnlyFs = afero.ReadOnlyFs

type ReadOnlyFile struct{}

// Sync implements afero.File.
func (ReadOnlyFile) Sync() error {
	return nil
}

// Truncate implements afero.File.
func (ReadOnlyFile) Truncate(int64) error {
	return syscall.EROFS
}

// Write implements afero.File.
func (ReadOnlyFile) Write([]byte) (n int, err error) {
	return 0, syscall.EROFS
}

// WriteAt implements afero.File.
func (ReadOnlyFile) WriteAt([]byte, int64) (n int, err error) {
	return 0, syscall.EROFS
}

// WriteString implements afero.File.
func (ReadOnlyFile) WriteString(string) (ret int, err error) {
	return 0, syscall.EROFS
}
