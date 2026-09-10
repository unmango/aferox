package internal_test

import (
	"os"

	filev1alpha1 "buf.build/gen/go/unmango/protofs/protocolbuffers/go/dev/unmango/file/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unmango/aferox/protofs/internal"
)

var _ = Describe("Util", func() {
	DescribeTable("OsFileMode",
		func(mode filev1alpha1.FileMode, expected os.FileMode) {
			actual := internal.OsFileMode(mode)

			Expect(actual).To(Equal(expected))
		},
		Entry(nil, filev1alpha1.FileMode_FILE_MODE_APPEND, os.ModeAppend),
		Entry(nil, filev1alpha1.FileMode_FILE_MODE_CHAR_DEVICE, os.ModeCharDevice),
		Entry(nil, filev1alpha1.FileMode_FILE_MODE_DEVICE, os.ModeDevice),
		Entry(nil, filev1alpha1.FileMode_FILE_MODE_DIR, os.ModeDir),
		Entry(nil, filev1alpha1.FileMode_FILE_MODE_EXCLUSIVE, os.ModeExclusive),
		Entry(nil, filev1alpha1.FileMode_FILE_MODE_IRREGULAR, os.ModeIrregular),
		Entry(nil, filev1alpha1.FileMode_FILE_MODE_NAMED_PIPE, os.ModeNamedPipe),
		Entry(nil, filev1alpha1.FileMode_FILE_MODE_PERM, os.ModePerm),
		Entry(nil, filev1alpha1.FileMode_FILE_MODE_SETGID, os.ModeSetgid),
		Entry(nil, filev1alpha1.FileMode_FILE_MODE_SETUID, os.ModeSetuid),
		Entry(nil, filev1alpha1.FileMode_FILE_MODE_SOCKET, os.ModeSocket),
		Entry(nil, filev1alpha1.FileMode_FILE_MODE_STICKY, os.ModeSticky),
		Entry(nil, filev1alpha1.FileMode_FILE_MODE_SYMLINK, os.ModeSymlink),
		Entry(nil, filev1alpha1.FileMode_FILE_MODE_TEMPORARY, os.ModeTemporary),
	)

	DescribeTable("ProtoFileMode",
		func(mode os.FileMode, expected filev1alpha1.FileMode) {
			actual := internal.ProtoFileMode(mode)

			Expect(actual).To(Equal(expected))
		},
		Entry(nil, os.ModeAppend, filev1alpha1.FileMode_FILE_MODE_APPEND),
		Entry(nil, os.ModeCharDevice, filev1alpha1.FileMode_FILE_MODE_CHAR_DEVICE),
		Entry(nil, os.ModeDevice, filev1alpha1.FileMode_FILE_MODE_DEVICE),
		Entry(nil, os.ModeDir, filev1alpha1.FileMode_FILE_MODE_DIR),
		Entry(nil, os.ModeExclusive, filev1alpha1.FileMode_FILE_MODE_EXCLUSIVE),
		Entry(nil, os.ModeIrregular, filev1alpha1.FileMode_FILE_MODE_IRREGULAR),
		Entry(nil, os.ModeNamedPipe, filev1alpha1.FileMode_FILE_MODE_NAMED_PIPE),
		Entry(nil, os.ModePerm, filev1alpha1.FileMode_FILE_MODE_PERM),
		Entry(nil, os.ModeSetgid, filev1alpha1.FileMode_FILE_MODE_SETGID),
		Entry(nil, os.ModeSetuid, filev1alpha1.FileMode_FILE_MODE_SETUID),
		Entry(nil, os.ModeSocket, filev1alpha1.FileMode_FILE_MODE_SOCKET),
		Entry(nil, os.ModeSticky, filev1alpha1.FileMode_FILE_MODE_STICKY),
		Entry(nil, os.ModeSymlink, filev1alpha1.FileMode_FILE_MODE_SYMLINK),
		Entry(nil, os.ModeTemporary, filev1alpha1.FileMode_FILE_MODE_TEMPORARY),
	)
})
