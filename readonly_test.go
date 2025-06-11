package aferox_test

import (
	"syscall"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/unmango/aferox"
)

var _ = Describe("Readonly", func() {
	var file *aferox.ReadOnlyFile

	BeforeEach(func() {
		file = &aferox.ReadOnlyFile{}
	})

	It("should error on readdir", func() {
		_, err := file.Readdir(0)

		Expect(err).To(MatchError(syscall.EPERM))
	})

	It("should error on readdir names", func() {
		_, err := file.Readdirnames(0)

		Expect(err).To(MatchError(syscall.EPERM))
	})

	It("should error on seek", func() {
		_, err := file.Seek(0, 0)

		Expect(err).To(MatchError(syscall.EPERM))
	})

	It("should error on sync", func() {
		Expect(file.Sync()).To(MatchError(syscall.EPERM))
	})

	It("should error on truncate", func() {
		Expect(file.Truncate(0)).To(MatchError(syscall.EPERM))
	})

	It("should error on write", func() {
		_, err := file.Write(nil)

		Expect(err).To(MatchError(syscall.EPERM))
	})

	It("should error on write at", func() {
		_, err := file.WriteAt(nil, 0)

		Expect(err).To(MatchError(syscall.EPERM))
	})

	It("should error on write string", func() {
		_, err := file.WriteString("")

		Expect(err).To(MatchError(syscall.EPERM))
	})
})
