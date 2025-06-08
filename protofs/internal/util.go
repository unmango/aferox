package internal

import (
	"os"

	fsv1alpha1 "buf.build/gen/go/unmango/protofs/protocolbuffers/go/dev/unmango/fs/v1alpha1"
)

func OsFileMode(mode fsv1alpha1.FileMode) os.FileMode {
	return os.FileMode(mode)
}

func ProtoFileMode(mode os.FileMode) fsv1alpha1.FileMode {
	return fsv1alpha1.FileMode(mode)
}
