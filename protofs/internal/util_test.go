package internal_test

import (
	"os"

	fsv1alpha1 "buf.build/gen/go/unmango/protofs/protocolbuffers/go/dev/unmango/fs/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unmango/aferox/protofs/internal"
)

var _ = Describe("Util", func() {
	DescribeTable("OsFileMode",
		func(mode fsv1alpha1.FileMode, expected os.FileMode) {
			actual := internal.OsFileMode(mode)

			Expect(actual).To(Equal(expected))
		},
		Entry(nil, fsv1alpha1.FileMode_FILE_MODE_APPEND, os.ModeAppend),
		Entry(nil, fsv1alpha1.FileMode_FILE_MODE_CHAR_DEVICE, os.ModeCharDevice),
		Entry(nil, fsv1alpha1.FileMode_FILE_MODE_DEVICE, os.ModeDevice),
		Entry(nil, fsv1alpha1.FileMode_FILE_MODE_DIR, os.ModeDir),
		Entry(nil, fsv1alpha1.FileMode_FILE_MODE_EXCLUSIVE, os.ModeExclusive),
		Entry(nil, fsv1alpha1.FileMode_FILE_MODE_IRREGULAR, os.ModeIrregular),
		Entry(nil, fsv1alpha1.FileMode_FILE_MODE_NAMED_PIPE, os.ModeNamedPipe),
		Entry(nil, fsv1alpha1.FileMode_FILE_MODE_PERM, os.ModePerm),
		Entry(nil, fsv1alpha1.FileMode_FILE_MODE_SETGID, os.ModeSetgid),
		Entry(nil, fsv1alpha1.FileMode_FILE_MODE_SETUID, os.ModeSetuid),
		Entry(nil, fsv1alpha1.FileMode_FILE_MODE_SOCKET, os.ModeSocket),
		Entry(nil, fsv1alpha1.FileMode_FILE_MODE_STICKY, os.ModeSticky),
		Entry(nil, fsv1alpha1.FileMode_FILE_MODE_SYMLINK, os.ModeSymlink),
		Entry(nil, fsv1alpha1.FileMode_FILE_MODE_TEMPORARY, os.ModeTemporary),
	)

	DescribeTable("ProtoFileMode",
		func(mode os.FileMode, expected fsv1alpha1.FileMode) {
			actual := internal.ProtoFileMode(mode)

			Expect(actual).To(Equal(expected))
		},
		Entry(nil, os.ModeAppend, fsv1alpha1.FileMode_FILE_MODE_APPEND),
		Entry(nil, os.ModeCharDevice, fsv1alpha1.FileMode_FILE_MODE_CHAR_DEVICE),
		Entry(nil, os.ModeDevice, fsv1alpha1.FileMode_FILE_MODE_DEVICE),
		Entry(nil, os.ModeDir, fsv1alpha1.FileMode_FILE_MODE_DIR),
		Entry(nil, os.ModeExclusive, fsv1alpha1.FileMode_FILE_MODE_EXCLUSIVE),
		Entry(nil, os.ModeIrregular, fsv1alpha1.FileMode_FILE_MODE_IRREGULAR),
		Entry(nil, os.ModeNamedPipe, fsv1alpha1.FileMode_FILE_MODE_NAMED_PIPE),
		Entry(nil, os.ModePerm, fsv1alpha1.FileMode_FILE_MODE_PERM),
		Entry(nil, os.ModeSetgid, fsv1alpha1.FileMode_FILE_MODE_SETGID),
		Entry(nil, os.ModeSetuid, fsv1alpha1.FileMode_FILE_MODE_SETUID),
		Entry(nil, os.ModeSocket, fsv1alpha1.FileMode_FILE_MODE_SOCKET),
		Entry(nil, os.ModeSticky, fsv1alpha1.FileMode_FILE_MODE_STICKY),
		Entry(nil, os.ModeSymlink, fsv1alpha1.FileMode_FILE_MODE_SYMLINK),
		Entry(nil, os.ModeTemporary, fsv1alpha1.FileMode_FILE_MODE_TEMPORARY),
	)
})
