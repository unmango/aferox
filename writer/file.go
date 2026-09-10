package writer

import (
	"io"
	"io/fs"
	"syscall"
)

type File struct {
	writer io.Writer
	name   string
}

// Close implements afero.File.
func (f *File) Close() error {
	return nil
}

// Name implements afero.File.
func (f *File) Name() string {
	return f.name
}

// Read implements afero.File.
func (f *File) Read([]byte) (n int, err error) {
	return 0, syscall.EROFS
}

// ReadAt implements afero.File.
func (f *File) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, syscall.EROFS
}

// Readdir implements afero.File.
func (f *File) Readdir(count int) ([]fs.FileInfo, error) {
	return nil, syscall.EROFS
}

// Readdirnames implements afero.File.
func (f *File) Readdirnames(n int) ([]string, error) {
	return nil, syscall.EROFS
}

// Seek implements afero.File.
func (f *File) Seek(offset int64, whence int) (int64, error) {
	return 0, syscall.EROFS
}

// Stat implements afero.File.
func (f *File) Stat() (fs.FileInfo, error) {
	return nil, syscall.EROFS
}

// Sync implements afero.File.
func (f *File) Sync() error {
	return syscall.EROFS
}

// Truncate implements afero.File.
func (f *File) Truncate(size int64) error {
	return syscall.EROFS
}

// Write implements afero.File.
func (f *File) Write(p []byte) (n int, err error) {
	return f.writer.Write(p)
}

// WriteAt implements afero.File.
func (f *File) WriteAt(p []byte, off int64) (n int, err error) {
	if wa, ok := f.writer.(io.WriterAt); ok {
		return wa.WriteAt(p, off)
	}

	return 0, syscall.EROFS
}

// WriteString implements afero.File.
func (f *File) WriteString(s string) (ret int, err error) {
	return io.WriteString(f.writer, s)
}
