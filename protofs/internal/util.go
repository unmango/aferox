package internal

import (
	"os"

	filev1alpha1 "buf.build/gen/go/unmango/protofs/protocolbuffers/go/dev/unmango/file/v1alpha1"
)

func OsFileMode(mode filev1alpha1.FileMode) os.FileMode {
	return os.FileMode(mode)
}

func ProtoFileMode(mode os.FileMode) filev1alpha1.FileMode {
	return filev1alpha1.FileMode(mode)
}
