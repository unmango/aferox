package image

import (
	"archive/tar"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/spf13/afero"
	"github.com/spf13/afero/tarfs"
	"github.com/unmango/aferox/containerregistry/v1/layer"
)

// The afero version of this upstream issue for fs.FS
// https://github.com/google/go-containerregistry/issues/921#issuecomment-769252935

func FromFs(fs afero.Fs) (v1.Image, error) {
	l, err := layer.FromFs(fs)
	if err != nil {
		return nil, err
	}

	return mutate.AppendLayers(empty.Image, l)
}

func ToFs(img v1.Image) (afero.Fs, error) {
	rc := mutate.Extract(img)
	defer rc.Close()

	return tarfs.New(tar.NewReader(rc)), nil
}
