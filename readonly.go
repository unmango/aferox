package aferox

import (
	"os"
	"syscall"
	"time"

	"github.com/spf13/afero"
)

// ReadOnlyFs can be embedded in structs to assist with defining read-only implementations of afero.Fs
type ReadOnlyFs struct{}

// Chmod implements afero.Fs.
func (r *ReadOnlyFs) Chmod(name string, mode os.FileMode) error {
	return syscall.EPERM
}

// Chown implements afero.Fs.
func (r *ReadOnlyFs) Chown(name string, uid int, gid int) error {
	return syscall.EPERM
}

// Chtimes implements afero.Fs.
func (r *ReadOnlyFs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return syscall.EPERM
}

// Create implements afero.Fs.
func (r *ReadOnlyFs) Create(name string) (afero.File, error) {
	return nil, syscall.EPERM
}

// Mkdir implements afero.Fs.
func (r *ReadOnlyFs) Mkdir(name string, perm os.FileMode) error {
	return syscall.EPERM
}

// MkdirAll implements afero.Fs.
func (r *ReadOnlyFs) MkdirAll(path string, perm os.FileMode) error {
	return syscall.EPERM
}

// Remove implements afero.Fs.
func (r *ReadOnlyFs) Remove(name string) error {
	return syscall.EPERM
}

// RemoveAll implements afero.Fs.
func (r *ReadOnlyFs) RemoveAll(path string) error {
	return syscall.EPERM
}

// Rename implements afero.Fs.
func (r *ReadOnlyFs) Rename(oldname string, newname string) error {
	return syscall.EPERM
}

// ReadOnlyFile can be embedded in structs to assist with defining read-only implementations of afero.File
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
