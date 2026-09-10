package internal_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"

	"github.com/unmango/aferox/internal"
	"github.com/unmango/aferox/testing/gfs"
)

var _ = Describe("Copy", func() {
	It("should error when src doesn't exist", func() {
		src := afero.NewBasePathFs(afero.NewMemMapFs(), "blah")
		dest := afero.NewMemMapFs()

		err := internal.Copy(src, dest)

		Expect(err).To(MatchError("open blah: file does not exist"))
	})

	It("should create dest when it doesn't exist", func() {
		src := afero.NewMemMapFs()
		err := afero.WriteFile(src, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		base := afero.NewMemMapFs()
		dest := afero.NewBasePathFs(base, "blah")

		err = internal.Copy(src, dest)

		Expect(err).NotTo(HaveOccurred())
		Expect(dest).To(gfs.ContainFile("test.txt"))
		Expect(base).To(gfs.ContainFile("blah/test.txt"))
	})

	It("should copy files", func() {
		src := afero.NewMemMapFs()
		err := afero.WriteFile(src, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		dest := afero.NewMemMapFs()

		err = internal.Copy(src, dest)

		Expect(err).NotTo(HaveOccurred())
		Expect(dest).To(gfs.ContainFileWithBytes("test.txt", []byte("testing")))
	})

	It("should copy directories", func() {
		src := afero.NewMemMapFs()
		err := src.Mkdir("test", os.ModeDir)
		Expect(err).NotTo(HaveOccurred())
		dest := afero.NewMemMapFs()

		err = internal.Copy(src, dest)

		Expect(err).NotTo(HaveOccurred())
		stat, err := dest.Stat("test")
		Expect(err).NotTo(HaveOccurred())
		Expect(stat.IsDir()).To(BeTrueBecause("the directory is created"))
	})

	It("should copy directories with files", func() {
		src := afero.NewMemMapFs()
		err := src.Mkdir("test", os.ModeDir)
		Expect(err).NotTo(HaveOccurred())
		err = afero.WriteFile(src, "test/test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		dest := afero.NewMemMapFs()

		err = internal.Copy(src, dest)

		Expect(err).NotTo(HaveOccurred())
		stat, err := dest.Stat("test")
		Expect(err).NotTo(HaveOccurred())
		Expect(stat.IsDir()).To(BeTrueBecause("the directory is created"))
		Expect(dest).To(gfs.ContainFileWithBytes("test/test.txt", []byte("testing")))
	})

	It("should copy multiple files", func() {
		src := afero.NewMemMapFs()
		err := afero.WriteFile(src, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		err = afero.WriteFile(src, "test2.txt", []byte("testing2"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		dest := afero.NewMemMapFs()

		err = internal.Copy(src, dest)

		Expect(err).NotTo(HaveOccurred())
		Expect(dest).To(gfs.ContainFileWithBytes("test.txt", []byte("testing")))
		Expect(dest).To(gfs.ContainFileWithBytes("test2.txt", []byte("testing2")))
	})

	It("should copy a directory structure", func() {
		src := afero.NewMemMapFs()
		err := afero.WriteFile(src, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		err = src.MkdirAll("test/other", os.ModeDir)
		Expect(err).NotTo(HaveOccurred())
		err = afero.WriteFile(src, "test/test2.txt", []byte("testing2"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		err = afero.WriteFile(src, "test/other/test3.txt", []byte("testing3"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		dest := afero.NewMemMapFs()

		err = internal.Copy(src, dest)

		Expect(err).NotTo(HaveOccurred())
		Expect(dest).To(gfs.ContainFileWithBytes("test.txt", []byte("testing")))
		stat, err := dest.Stat("test")
		Expect(err).NotTo(HaveOccurred())
		Expect(stat.IsDir()).To(BeTrueBecause("the first directory is created"))
		Expect(dest).To(gfs.ContainFileWithBytes("test/test2.txt", []byte("testing2")))
		stat, err = dest.Stat("test/other")
		Expect(err).NotTo(HaveOccurred())
		Expect(stat.IsDir()).To(BeTrueBecause("the second directory is created"))
		Expect(dest).To(gfs.ContainFileWithBytes("test/other/test3.txt", []byte("testing3")))
	})
})
