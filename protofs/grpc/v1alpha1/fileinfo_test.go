package protofsv1alpha1_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	filev1alpha1 "buf.build/gen/go/unmango/protofs/protocolbuffers/go/dev/unmango/file/v1alpha1"
	protofsv1alpha1 "github.com/unmango/aferox/protofs/grpc/v1alpha1"
	"github.com/unmango/aferox/protofs/internal"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ = Describe("FileInfo", func() {
	It("should map values", func() {
		proto := &filev1alpha1.FileInfo{
			Name:    "test-name",
			Size:    69,
			Mode:    filev1alpha1.FileMode_FILE_MODE_APPEND,
			ModTime: timestamppb.Now(),
			IsDir:   true,
		}

		fi := &protofsv1alpha1.FileInfo{Proto: proto}

		Expect(fi.Name()).To(Equal(proto.Name))
		Expect(fi.Size()).To(Equal(proto.Size))
		Expect(fi.Mode()).To(Equal(internal.OsFileMode(proto.Mode)))
		Expect(fi.ModTime()).To(Equal(proto.ModTime.AsTime()))
		Expect(fi.IsDir()).To(Equal(proto.IsDir))
	})
})
