package aferox_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"

	"github.com/unmango/aferox"
)

var _ = Describe("Single", func() {
	Describe("StatSingle", func() {
		It("should stat an Fs with a single file", func() {
			fsys := afero.NewMemMapFs()
			_, err := fsys.Create("test.txt")
			Expect(err).NotTo(HaveOccurred())

			info, err := aferox.StatFirst(fsys, "")

			Expect(err).NotTo(HaveOccurred())
			Expect(info.Name()).To(Equal("test.txt"))
		})

		It("should stat an Fs with a single directory", func() {
			fsys := afero.NewMemMapFs()
			err := fsys.Mkdir("test", os.ModeDir)
			Expect(err).NotTo(HaveOccurred())

			info, err := aferox.StatFirst(fsys, "")

			Expect(err).NotTo(HaveOccurred())
			Expect(info.Name()).To(Equal("test"))
		})

		It("should not error when Fs contains multiple files", func() {
			fsys := afero.NewMemMapFs()
			_, err := fsys.Create("test.txt")
			Expect(err).NotTo(HaveOccurred())
			_, err = fsys.Create("oops.txt")
			Expect(err).NotTo(HaveOccurred())

			_, err = aferox.StatFirst(fsys, "")

			Expect(err).NotTo(HaveOccurred())
		})

		It("should error when Fs contains no files", func() {
			fsys := afero.NewMemMapFs()

			_, err := aferox.StatFirst(fsys, "")

			Expect(err).To(HaveOccurred())
		})

		When("SkipDirs is provided", func() {
			It("should stat the first file", func() {
				fsys := afero.NewMemMapFs()
				err := fsys.Mkdir("test", os.ModeDir)
				Expect(err).NotTo(HaveOccurred())
				_, err = fsys.Create("test/test.txt")

				info, err := aferox.StatFirst(fsys, "", aferox.SkipDirs)

				Expect(err).NotTo(HaveOccurred())
				Expect(info.Name()).To(Equal("test.txt"))
			})

			It("should error when only directories exist", func() {
				fsys := afero.NewMemMapFs()
				err := fsys.Mkdir("test", os.ModeDir)

				_, err = aferox.StatFirst(fsys, "", aferox.SkipDirs)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("OpenSingle", func() {
		It("should open in an Fs with a single file", func() {
			fsys := afero.NewMemMapFs()
			_, err := fsys.Create("test.txt")
			Expect(err).NotTo(HaveOccurred())

			info, err := aferox.OpenFirst(fsys, "")

			Expect(err).NotTo(HaveOccurred())
			Expect(info.Name()).To(Equal("test.txt"))
		})

		It("should open in an Fs with a single directory", func() {
			fsys := afero.NewMemMapFs()
			err := fsys.Mkdir("test", os.ModeDir)
			Expect(err).NotTo(HaveOccurred())

			info, err := aferox.OpenFirst(fsys, "")

			Expect(err).NotTo(HaveOccurred())
			Expect(info.Name()).To(Equal("test"))
		})

		It("should not error when Fs contains multiple files", func() {
			fsys := afero.NewMemMapFs()
			_, err := fsys.Create("test.txt")
			Expect(err).NotTo(HaveOccurred())
			_, err = fsys.Create("oops.txt")
			Expect(err).NotTo(HaveOccurred())

			_, err = aferox.OpenFirst(fsys, "")

			Expect(err).NotTo(HaveOccurred())
		})

		It("should error when Fs contains no files", func() {
			fsys := afero.NewMemMapFs()

			_, err := aferox.OpenFirst(fsys, "")

			Expect(err).To(HaveOccurred())
		})

		When("SkipDirs is provided", func() {
			It("should stat the first file", func() {
				fsys := afero.NewMemMapFs()
				err := fsys.Mkdir("test", os.ModeDir)
				Expect(err).NotTo(HaveOccurred())
				_, err = fsys.Create("test/test.txt")

				info, err := aferox.OpenFirst(fsys, "", aferox.SkipDirs)

				Expect(err).NotTo(HaveOccurred())
				Expect(info.Name()).To(Equal("test/test.txt"))
			})

			It("should error when only directories exist", func() {
				fsys := afero.NewMemMapFs()
				err := fsys.Mkdir("test", os.ModeDir)

				_, err = aferox.OpenFirst(fsys, "", aferox.SkipDirs)

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
