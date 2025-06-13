package protofsv1alpha1

import (
	"io/fs"
	"time"

	filev1alpha1 "buf.build/gen/go/unmango/protofs/protocolbuffers/go/dev/unmango/file/v1alpha1"
	"github.com/unmango/aferox/protofs/internal"
)

type FileInfo struct {
	Proto *filev1alpha1.FileInfo
}

// IsDir implements fs.FileInfo.
func (f FileInfo) IsDir() bool {
	return f.Proto.IsDir
}

// ModTime implements fs.FileInfo.
func (f FileInfo) ModTime() time.Time {
	return f.Proto.ModTime.AsTime()
}

// Mode implements fs.FileInfo.
func (f FileInfo) Mode() fs.FileMode {
	return internal.OsFileMode(f.Proto.Mode)
}

// Name implements fs.FileInfo.
func (f FileInfo) Name() string {
	return f.Proto.Name
}

// Size implements fs.FileInfo.
func (f FileInfo) Size() int64 {
	return f.Proto.Size
}

// Sys implements fs.FileInfo.
func (f FileInfo) Sys() any {
	return f.Proto.Sys
}
