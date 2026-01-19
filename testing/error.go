package testing

import (
	"io/fs"
	"os"
	"time"

	"github.com/spf13/afero"
	"github.com/spf13/afero/mem"
)

type ErrorFs struct {
	afero.Fs

	ChmodErr     error
	ChownErr     error
	ChtimesErr   error
	CreateErr    error
	MkdirErr     error
	MkdirAllErr  error
	OpenErr      error
	OpenFileErr  error
	RemoveErr    error
	RemoveAllErr error
	RenameErr    error
	StatErr      error
}

// Chmod implements afero.Fs.
func (f *ErrorFs) Chmod(name string, mode fs.FileMode) error {
	if f.ChmodErr != nil {
		return f.ChmodErr
	}
	return f.fs().Chmod(name, mode)
}

// Chown implements afero.Fs.
func (f *ErrorFs) Chown(name string, uid int, gid int) error {
	if f.ChownErr != nil {
		return f.ChownErr
	}
	return f.fs().Chown(name, uid, gid)
}

// Chtimes implements afero.Fs.
func (f *ErrorFs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	if f.ChtimesErr != nil {
		return f.ChtimesErr
	}
	return f.fs().Chtimes(name, atime, mtime)
}

// Create implements afero.Fs.
func (f *ErrorFs) Create(name string) (afero.File, error) {
	if f.CreateErr != nil {
		return nil, f.CreateErr
	}
	return f.fs().Create(name)
}

// Mkdir implements afero.Fs.
func (f *ErrorFs) Mkdir(name string, perm fs.FileMode) error {
	if f.MkdirErr != nil {
		return f.MkdirErr
	}
	return f.fs().Mkdir(name, perm)
}

// MkdirAll implements afero.Fs.
func (f *ErrorFs) MkdirAll(path string, perm fs.FileMode) error {
	if f.MkdirAllErr != nil {
		return f.MkdirAllErr
	}
	return f.fs().MkdirAll(path, perm)
}

// Name implements afero.Fs.
func (f *ErrorFs) Name() string {
	return "ErrorFs"
}

// Open implements afero.Fs.
func (f *ErrorFs) Open(name string) (afero.File, error) {
	if f.OpenErr != nil {
		return nil, f.OpenErr
	}
	return f.fs().Open(name)
}

// OpenFile implements afero.Fs.
func (f *ErrorFs) OpenFile(name string, flag int, perm fs.FileMode) (afero.File, error) {
	if f.OpenFileErr != nil {
		return nil, f.OpenFileErr
	}
	return f.fs().OpenFile(name, flag, perm)
}

// Remove implements afero.Fs.
func (f *ErrorFs) Remove(name string) error {
	if f.RemoveErr != nil {
		return f.RemoveErr
	}
	return f.fs().Remove(name)
}

// RemoveAll implements afero.Fs.
func (f *ErrorFs) RemoveAll(path string) error {
	if f.RemoveAllErr != nil {
		return f.RemoveAllErr
	}
	return f.fs().RemoveAll(path)
}

// Rename implements afero.Fs.
func (f *ErrorFs) Rename(oldname string, newname string) error {
	if f.RenameErr != nil {
		return f.RenameErr
	}
	return f.fs().Rename(oldname, newname)
}

// Stat implements afero.Fs.
func (f *ErrorFs) Stat(name string) (fs.FileInfo, error) {
	if f.StatErr != nil {
		return nil, f.StatErr
	}
	return f.fs().Stat(name)
}

func (f *ErrorFs) fs() afero.Fs {
	if f.Fs == nil {
		f.Fs = afero.NewMemMapFs()
	}
	return f.Fs
}

type ErrorFile struct {
	afero.File

	CloseErr        error
	ReadErr         error
	ReadAtErr       error
	ReaddirErr      error
	ReaddirnamesErr error
	SeekErr         error
	StatErr         error
	SyncErr         error
	TruncateErr     error
	WriteErr        error
	WriteAtErr      error
	WriteStringErr  error
}

// Close implements [afero.File].
func (e *ErrorFile) Close() error {
	if e.CloseErr != nil {
		return e.CloseErr
	}
	return e.file().Close()
}

// Read implements [afero.File].
func (e *ErrorFile) Read(p []byte) (n int, err error) {
	if e.ReadErr != nil {
		return 0, e.ReadErr
	}
	return e.file().Read(p)
}

// ReadAt implements [afero.File].
func (e *ErrorFile) ReadAt(p []byte, off int64) (n int, err error) {
	if e.ReadAtErr != nil {
		return 0, e.ReadAtErr
	}
	return e.file().ReadAt(p, off)
}

// Readdir implements [afero.File].
func (e *ErrorFile) Readdir(count int) ([]os.FileInfo, error) {
	if e.ReaddirErr != nil {
		return nil, e.ReaddirErr
	}
	return e.file().Readdir(count)
}

// Readdirnames implements [afero.File].
func (e *ErrorFile) Readdirnames(n int) ([]string, error) {
	if e.ReaddirnamesErr != nil {
		return nil, e.ReaddirnamesErr
	}
	return e.file().Readdirnames(n)
}

// Seek implements [afero.File].
func (e *ErrorFile) Seek(offset int64, whence int) (int64, error) {
	if e.SeekErr != nil {
		return 0, e.SeekErr
	}
	return e.file().Seek(offset, whence)
}

// Stat implements [afero.File].
func (e *ErrorFile) Stat() (os.FileInfo, error) {
	if e.StatErr != nil {
		return nil, e.StatErr
	}
	return e.file().Stat()
}

// Sync implements [afero.File].
func (e *ErrorFile) Sync() error {
	if e.SyncErr != nil {
		return e.SyncErr
	}
	return e.file().Sync()
}

// Truncate implements [afero.File].
func (e *ErrorFile) Truncate(size int64) error {
	if e.TruncateErr != nil {
		return e.TruncateErr
	}
	return e.file().Truncate(size)
}

// Write implements [afero.File].
func (e *ErrorFile) Write(p []byte) (n int, err error) {
	if e.WriteErr != nil {
		return 0, e.WriteErr
	}
	return e.file().Write(p)
}

// WriteAt implements [afero.File].
func (e *ErrorFile) WriteAt(p []byte, off int64) (n int, err error) {
	if e.WriteAtErr != nil {
		return 0, e.WriteAtErr
	}
	return e.file().WriteAt(p, off)
}

// WriteString implements [afero.File].
func (e *ErrorFile) WriteString(s string) (ret int, err error) {
	if e.WriteStringErr != nil {
		return 0, e.WriteStringErr
	}
	return e.file().WriteString(s)
}

func (e *ErrorFile) file() afero.File {
	if e.File == nil {
		e.File = mem.NewFileHandle(mem.CreateFile("ErrorFs"))
	}
	return e.File
}
