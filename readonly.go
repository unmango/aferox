package aferox

import (
	"os"
	"syscall"
)

// ReadOnlyFile can be embedded to assist with defining read-only implementations of afero.File
type ReadOnlyFile struct{}

// Readdir implements afero.File.
func (s *ReadOnlyFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, syscall.EPERM
}

// Readdirnames implements afero.File.
func (s *ReadOnlyFile) Readdirnames(n int) ([]string, error) {
	return nil, syscall.EPERM
}

// Seek implements afero.File.
func (r *ReadOnlyFile) Seek(offset int64, whence int) (int64, error) {
	return 0, syscall.EPERM
}

// Sync implements afero.File.
func (r *ReadOnlyFile) Sync() error {
	return syscall.EPERM
}

// Truncate implements afero.File.
func (r *ReadOnlyFile) Truncate(size int64) error {
	return syscall.EPERM
}

// Write implements afero.File.
func (r *ReadOnlyFile) Write(p []byte) (n int, err error) {
	return 0, syscall.EPERM
}

// WriteAt implements afero.File.
func (r *ReadOnlyFile) WriteAt(p []byte, off int64) (n int, err error) {
	return 0, syscall.EPERM
}

// WriteString implements afero.File.
func (r *ReadOnlyFile) WriteString(s string) (ret int, err error) {
	return 0, syscall.EPERM
}
