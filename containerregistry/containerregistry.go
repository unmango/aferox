package containerregistry

import (
	"github.com/unmango/aferox/containerregistry/v1/image"
	"github.com/unmango/aferox/containerregistry/v1/layer"
)

var (
	ImageToFs   = image.ToFs
	LayerFromFs = layer.FromFs
	LayerToFs   = layer.ToFs
)
