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

func NewWriter(w io.Writer) *tar.Writer {
	return tar.NewWriter(w)
}

func NewReader(r io.Reader) *tar.Reader {
	return tar.NewReader(r)
}
