package tar

import (
	"archive/tar"
	"io"
)

const (
	TypeReg           = tar.TypeReg
	TypeLink          = tar.TypeLink
	TypeSymlink       = tar.TypeSymlink
	TypeChar          = tar.TypeChar
	TypeBlock         = tar.TypeBlock
	TypeDir           = tar.TypeDir
	TypeFifo          = tar.TypeFifo
	TypeCont          = tar.TypeCont
	TypeXHeader       = tar.TypeXHeader
	TypeXGlobalHeader = tar.TypeXGlobalHeader
	TypeGNUSparse     = tar.TypeGNUSparse
	TypeGNULongName   = tar.TypeGNULongName
	TypeGNULongLink   = tar.TypeGNULongLink
)

// NewWriter returns a new tar.Writer that writes a tar archive to w.
//
// This function is a thin wrapper around archive/tar.NewWriter and exists to
// provide a stable, package-local API for tar output. Callers should prefer
// using this constructor instead of archive/tar.NewWriter directly so that
// tar handling can be centralized in this package and extended in the future
// without requiring changes at call sites.
func NewWriter(w io.Writer) *tar.Writer {
	return tar.NewWriter(w)
}

// NewReader returns a new tar.Reader that reads a tar archive from r.
//
// This function is a thin wrapper around archive/tar.NewReader and is
// provided to mirror NewWriter and keep tar I/O behind this package's API.
// Using this constructor instead of archive/tar.NewReader directly helps keep
// callers decoupled from the underlying tar implementation.
func NewReader(r io.Reader) *tar.Reader {
	return tar.NewReader(r)
}
