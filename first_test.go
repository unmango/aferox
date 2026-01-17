package aferox_test

import (
	"errors"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/spf13/afero"
	"github.com/unmango/aferox"
)

// failingOpenFs wraps an afero.Fs and makes Open fail for non-root paths
type failingOpenFs struct{ afero.Fs }

func (f *failingOpenFs) Open(name string) (afero.File, error) {
	if name != "" && name != "." {
		return nil, errors.New("open failed")
	}
	return f.Fs.Open(name)
}

var _ = Describe("First", func() {
	Describe("StatFirst", func() {
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

		It("should return error when walking non-existent path", func() {
			fs := afero.NewReadOnlyFs(afero.NewMemMapFs())

			_, err := aferox.StatFirst(fs, "nonexistent")

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

	Describe("OpenFirst", func() {
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

		It("should return error when walking non-existent path", func() {
			fs := afero.NewReadOnlyFs(afero.NewMemMapFs())

			_, err := aferox.OpenFirst(fs, "nonexistent")

			Expect(err).To(HaveOccurred())
		})

		It("should handle Open errors in OpenFirst", func() {
			baseFs := afero.NewMemMapFs()
			_, err := baseFs.Create("test.txt")
			Expect(err).NotTo(HaveOccurred())

			fs := &failingOpenFs{Fs: baseFs}

			_, err = aferox.OpenFirst(fs, "")

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("open failed"))
		})
	})
})
