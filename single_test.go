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

			info, err := aferox.StatSingle(fsys, "")

			Expect(err).NotTo(HaveOccurred())
			Expect(info.Name()).To(Equal("test.txt"))
		})

		It("should stat an Fs with a single directory", func() {
			fsys := afero.NewMemMapFs()
			err := fsys.Mkdir("test", os.ModeDir)
			Expect(err).NotTo(HaveOccurred())

			info, err := aferox.StatSingle(fsys, "")

			Expect(err).NotTo(HaveOccurred())
			Expect(info.Name()).To(Equal("test"))
		})

		It("should error when Fs contains multiple files", func() {
			fsys := afero.NewMemMapFs()
			_, err := fsys.Create("test.txt")
			Expect(err).NotTo(HaveOccurred())
			_, err = fsys.Create("oops.txt")
			Expect(err).NotTo(HaveOccurred())

			_, err = aferox.StatSingle(fsys, "")

			Expect(err).To(HaveOccurred())
		})

		It("should error when Fs contains no files", func() {
			fsys := afero.NewMemMapFs()

			_, err := aferox.StatSingle(fsys, "")

			Expect(err).To(HaveOccurred())
		})

		When("SkipDirs is provided", func() {
			It("should stat the first file", func() {
				fsys := afero.NewMemMapFs()
				err := fsys.Mkdir("test", os.ModeDir)
				Expect(err).NotTo(HaveOccurred())
				_, err = fsys.Create("test/test.txt")

				info, err := aferox.StatSingle(fsys, "", aferox.SkipDirs)

				Expect(err).NotTo(HaveOccurred())
				Expect(info.Name()).To(Equal("test.txt"))
			})

			It("should error when only directories exist", func() {
				fsys := afero.NewMemMapFs()
				err := fsys.Mkdir("test", os.ModeDir)

				_, err = aferox.StatSingle(fsys, "", aferox.SkipDirs)

				Expect(err).To(HaveOccurred())
			})
		})

		It("should return error when walking non-existent path", func() {
			fs := afero.NewReadOnlyFs(afero.NewMemMapFs())

			_, err := aferox.StatSingle(fs, "nonexistent")

			Expect(err).To(HaveOccurred())
		})
	})

	Describe("OpenSingle", func() {
		It("should open in an Fs with a single file", func() {
			fsys := afero.NewMemMapFs()
			_, err := fsys.Create("test.txt")
			Expect(err).NotTo(HaveOccurred())

			info, err := aferox.OpenSingle(fsys, "")

			Expect(err).NotTo(HaveOccurred())
			Expect(info.Name()).To(Equal("test.txt"))
		})

		It("should open in an Fs with a single directory", func() {
			fsys := afero.NewMemMapFs()
			err := fsys.Mkdir("test", os.ModeDir)
			Expect(err).NotTo(HaveOccurred())

			info, err := aferox.OpenSingle(fsys, "")

			Expect(err).NotTo(HaveOccurred())
			Expect(info.Name()).To(Equal("test"))
		})

		It("should error when Fs contains multiple files", func() {
			fsys := afero.NewMemMapFs()
			_, err := fsys.Create("test.txt")
			Expect(err).NotTo(HaveOccurred())
			_, err = fsys.Create("oops.txt")
			Expect(err).NotTo(HaveOccurred())

			_, err = aferox.OpenSingle(fsys, "")

			Expect(err).To(HaveOccurred())
		})

		It("should error when Fs contains no files", func() {
			fsys := afero.NewMemMapFs()

			_, err := aferox.OpenSingle(fsys, "")

			Expect(err).To(HaveOccurred())
		})

		When("SkipDirs is provided", func() {
			It("should stat the first file", func() {
				fsys := afero.NewMemMapFs()
				err := fsys.Mkdir("test", os.ModeDir)
				Expect(err).NotTo(HaveOccurred())
				_, err = fsys.Create("test/test.txt")

				info, err := aferox.OpenSingle(fsys, "", aferox.SkipDirs)

				Expect(err).NotTo(HaveOccurred())
				Expect(info.Name()).To(Equal("test/test.txt"))
			})

			It("should error when only directories exist", func() {
				fsys := afero.NewMemMapFs()
				err := fsys.Mkdir("test", os.ModeDir)

				_, err = aferox.OpenSingle(fsys, "", aferox.SkipDirs)

				Expect(err).To(HaveOccurred())
			})
		})

		It("should return error when walking non-existent path", func() {
			fs := afero.NewReadOnlyFs(afero.NewMemMapFs())

			_, err := aferox.OpenSingle(fs, "nonexistent")

			Expect(err).To(HaveOccurred())
		})
	})

	Describe("errSingle", func() {
		It("should format error message correctly", func() {
			fsys := afero.NewMemMapFs()
			_, err := fsys.Create("file1.txt")
			Expect(err).NotTo(HaveOccurred())
			_, err = fsys.Create("file2.txt")
			Expect(err).NotTo(HaveOccurred())

			_, err = aferox.StatSingle(fsys, "")

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fs contains more than one entry"))
			Expect(err.Error()).To(ContainSubstring("had:"))
			Expect(err.Error()).To(ContainSubstring("found:"))
		})
	})
})
