package tar

import (
	"archive/tar"
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/afero"
)

// Fs implements [afero.Fs] as a write-only filesystem backed by an archive/tar.Writer.
//
// It is intended for creating tar archive entries via standard filesystem-style calls
// such as Create and OpenFile. Read, stat-based modification, and mutating operations
// (for example Chmod, Chown, Chtimes, Remove, RemoveAll, and Rename) are not supported
// on the underlying tar stream and will return appropriate permission or read-only
// filesystem errors (for example syscall.EPERM or syscall.EROFS).
//
// Callers should treat Fs as an append-only view of a tar archive and should not expect
// to be able to read back or modify previously written entries through this interface.
type Fs struct {
	m *sync.Mutex
	w *tar.Writer
}

// Chmod implements [afero.Fs].
func (f *Fs) Chmod(name string, mode os.FileMode) error {
	return syscall.EPERM
}

// Chown implements [afero.Fs].
func (f *Fs) Chown(name string, uid int, gid int) error {
	return syscall.EPERM
}

// Chtimes implements [afero.Fs].
func (f *Fs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return syscall.EPERM
}

// Create implements [afero.Fs].
func (f *Fs) Create(name string) (afero.File, error) {
	return newFile(name, f.w, f.m, 0644), nil
}

// Mkdir implements [afero.Fs].
func (f *Fs) Mkdir(name string, perm os.FileMode) error {
	f.m.Lock()
	defer f.m.Unlock()

	if err := f.w.WriteHeader(&tar.Header{
		Typeflag: tar.TypeDir,
		Name:     name,
		Mode:     int64(perm),
	}); err != nil {
		return fmt.Errorf("writing header: %w", err)
	}

	return nil
}

// MkdirAll implements [afero.Fs].
func (f *Fs) MkdirAll(path string, perm os.FileMode) error {
	return f.Mkdir(path, perm)
}

// Name implements [afero.Fs].
func (f *Fs) Name() string {
	return "tar.Writer"
}

// Open implements [afero.Fs].
func (f *Fs) Open(name string) (afero.File, error) {
	// This filesystem is write-only (backed by tar.Writer), so opening files
	// for reading is not supported.
	return nil, syscall.EROFS
}

// OpenFile implements [afero.Fs].
func (f *Fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	return newFile(name, f.w, f.m, perm), nil
}

// Remove implements [afero.Fs].
func (f *Fs) Remove(name string) error {
	return syscall.EROFS
}

// RemoveAll implements [afero.Fs].
func (f *Fs) RemoveAll(path string) error {
	return syscall.EROFS
}

// Rename implements [afero.Fs].
func (f *Fs) Rename(oldname string, newname string) error {
	return syscall.EROFS
}

// Stat implements [afero.Fs].
func (f *Fs) Stat(name string) (os.FileInfo, error) {
	return &FileInfo{name: name, w: f.w}, nil
}

// NewFs returns an [afero.Fs] implementation that writes its contents to the
// provided [tar.Writer]. The returned filesystem is write-only and exposes a
// subset of operations that add files and directories to the underlying tar
// stream.
//
// The caller retains ownership of the tar.Writer: NewFs does not close or
// flush the writer, and it does not perform any synchronization. The caller
// is responsible for calling Close on the tar.Writer when all filesystem
// operations are complete.
//
// Example:
//
//   tw := tar.NewWriter(dst)
//   defer tw.Close()
//
//   fs := NewFs(tw)
//   // Use any afero helpers with fs, for example:
//   //   afero.WriteFile(fs, "path/to/file.txt", []byte("data"), 0o644)
func NewFs(w *tar.Writer) afero.Fs {
	return &Fs{
		m: &sync.Mutex{},
		w: w,
	}
}
