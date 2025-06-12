package aferox_test

import (
	"os"
	"syscall"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/unmango/aferox"
)

var _ = Describe("Readonly", func() {
	Describe("ReadonlyFs", func() {
		var fs *aferox.ReadOnlyFs

		BeforeEach(func() {
			fs = &aferox.ReadOnlyFs{}
		})

		It("should error on Create", func() {
			Expect(fs.Create("")).To(MatchError(syscall.EPERM))
		})

		It("should error on Chmod", func() {
			Expect(fs.Chmod("", 0)).To(MatchError(syscall.EPERM))
		})

		It("should error on Chown", func() {
			Expect(fs.Chown("", 0, 0)).To(MatchError(syscall.EPERM))
		})

		It("should error on Chtimes", func() {
			Expect(fs.Chtimes("", time.Now(), time.Now())).To(MatchError(syscall.EPERM))
		})

		It("should error on Mkdir", func() {
			Expect(fs.Mkdir("", os.ModePerm)).To(MatchError(syscall.EPERM))
		})

		It("should error on MkdirAll", func() {
			Expect(fs.MkdirAll("", os.ModePerm)).To(MatchError(syscall.EPERM))
		})

		It("should error on Remove", func() {
			Expect(fs.Remove("")).To(MatchError(syscall.EPERM))
		})

		It("should error on RemoveAll", func() {
			Expect(fs.RemoveAll("")).To(MatchError(syscall.EPERM))
		})

		It("should error on Rename", func() {
			Expect(fs.Rename("", "")).To(MatchError(syscall.EPERM))
		})
	})

	Describe("ReadonlyFile", func() {
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
})
