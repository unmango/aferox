package tar

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"syscall"
)

type File struct {
	name string
	w    *tar.Writer
	buf  *bytes.Buffer
	perm os.FileMode

	once *sync.Once
}

func NewFile(name string, w *tar.Writer, perm os.FileMode) *File {
	return &File{
		name: name,
		w:    w,
		buf:  &bytes.Buffer{},
		perm: perm,
		once: &sync.Once{},
	}
}

func (f *File) flush() error {
	if err := f.w.WriteHeader(&tar.Header{
		Name:     f.name,
		Size:     int64(f.buf.Len()),
		Mode:     int64(f.perm),
		Typeflag: tar.TypeReg,
	}); err != nil {
		return fmt.Errorf("writing header: %w", err)
	}

	if _, err := io.Copy(f.w, f.buf); err != nil {
		return fmt.Errorf("copying file buffer: %w", err)
	}

	return nil
}

// Close implements [afero.File].
func (f *File) Close() error {
	var err error
	f.once.Do(func() {
		err = f.flush()
	})

	return err
}

// Name implements [afero.File].
func (f *File) Name() string {
	return f.name
}

// Read implements [afero.File].
func (f *File) Read(p []byte) (n int, err error) {
	return 0, syscall.EROFS
}

// ReadAt implements [afero.File].
func (f *File) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, syscall.EROFS
}

// Readdir implements [afero.File].
func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	return nil, syscall.EROFS
}

// Readdirnames implements [afero.File].
func (f *File) Readdirnames(n int) ([]string, error) {
	return nil, syscall.EROFS
}

// Seek implements [afero.File].
func (f *File) Seek(offset int64, whence int) (int64, error) {
	return 0, syscall.EROFS
}

// Stat implements [afero.File].
func (f *File) Stat() (os.FileInfo, error) {
	return &FileInfo{name: f.name, w: f.w}, nil
}

// Sync implements [afero.File].
func (f *File) Sync() error {
	return syscall.EROFS
}

// Truncate implements [afero.File].
func (f *File) Truncate(size int64) error {
	return syscall.EROFS
}

// Write implements [afero.File].
func (f *File) Write(p []byte) (n int, err error) {
	return f.buf.Write(p)
}

// WriteAt implements [afero.File].
func (f *File) WriteAt(p []byte, off int64) (n int, err error) {
	return 0, syscall.EROFS
}

// WriteString implements [afero.File].
func (f *File) WriteString(s string) (ret int, err error) {
	return io.WriteString(f.buf, s)
}
