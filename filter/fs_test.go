package filter_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"

	"github.com/unmango/aferox/filter"
)

var _ = Describe("Fs", func() {
	It("should not allow chmod-ing a filtered file", func() {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		filtered := filter.NewFs(fs, func(name string) bool {
			return name == "not-test.txt"
		})

		err = filtered.Chmod("test.txt", os.ModeAppend)

		Expect(err).To(MatchError(syscall.ENOENT))
	})

	It("should allow chmod-ing a non-filtered file", func() {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		filtered := filter.NewFs(fs, func(name string) bool {
			return name == "test.txt"
		})

		err = filtered.Chmod("test.txt", os.ModeAppend)

		Expect(err).NotTo(HaveOccurred())
	})

	It("should not allow chown-ing a filtered file", func() {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		filtered := filter.NewFs(fs, func(name string) bool {
			return name == "not-test.txt"
		})

		err = filtered.Chown("test.txt", 1001, 1001)

		Expect(err).To(MatchError(syscall.ENOENT))
	})

	It("should allow chown-ing a non-filtered file", func() {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		filtered := filter.NewFs(fs, func(name string) bool {
			return name == "test.txt"
		})

		err = filtered.Chown("test.txt", 1001, 1001)

		Expect(err).NotTo(HaveOccurred())
	})

	It("should not allow chtimes-ing a filtered file", func() {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		filtered := filter.NewFs(fs, func(name string) bool {
			return name == "not-test.txt"
		})

		err = filtered.Chtimes("test.txt", time.Now(), time.Now())

		Expect(err).To(MatchError(syscall.ENOENT))
	})

	It("should allow chtimes-ing a non-filtered file", func() {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		filtered := filter.NewFs(fs, func(name string) bool {
			return name == "test.txt"
		})

		err = filtered.Chtimes("test.txt", time.Now(), time.Now())

		Expect(err).NotTo(HaveOccurred())
	})

	It("should not allow creating a filtered file", func() {
		fs := afero.NewMemMapFs()
		filtered := filter.NewFs(fs, func(name string) bool {
			return name == "not-test.txt"
		})

		_, err := filtered.Create("test.txt")

		Expect(err).To(MatchError(syscall.ENOENT))
	})

	It("should allow creating a non-filtered file", func() {
		fs := afero.NewMemMapFs()
		filtered := filter.NewFs(fs, func(name string) bool {
			return name == "test.txt"
		})

		_, err := filtered.Create("test.txt")

		Expect(err).NotTo(HaveOccurred())
	})

	It("should include the source filesystem name", func() {
		fs := filter.NewFs(afero.NewMemMapFs(), func(s string) bool {
			return true
		})

		Expect(fs.Name()).To(Equal("Filter: MemMapFS"))
	})

	It("should not allow opening a filtered file", func() {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		filtered := filter.NewFs(fs, func(name string) bool {
			return name == "not-test.txt"
		})

		_, err = filtered.Open("test.txt")

		Expect(err).To(MatchError(syscall.ENOENT))
	})

	It("should allow opening a non-filtered file", func() {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		filtered := filter.NewFs(fs, func(name string) bool {
			return name == "test.txt"
		})

		_, err = filtered.Open("test.txt")

		Expect(err).NotTo(HaveOccurred())
	})

	It("should not allow open-file-ing a filtered file", func() {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		filtered := filter.NewFs(fs, func(name string) bool {
			return name == "not-test.txt"
		})

		_, err = filtered.OpenFile("test.txt", 69, os.ModeAppend)

		Expect(err).To(MatchError(syscall.ENOENT))
	})

	It("should allow open-file-ing a non-filtered file", func() {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		filtered := filter.NewFs(fs, func(name string) bool {
			return name == "test.txt"
		})

		_, err = filtered.OpenFile("test.txt", 69, os.ModeAppend)

		Expect(err).NotTo(HaveOccurred())
	})

	It("should not allow removing a filtered file", func() {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		filtered := filter.NewFs(fs, func(name string) bool {
			return name == "not-test.txt"
		})

		err = filtered.Remove("test.txt")

		Expect(err).To(MatchError(syscall.ENOENT))
	})

	It("should allow removing a non-filtered file", func() {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		filtered := filter.NewFs(fs, func(name string) bool {
			return name == "test.txt"
		})

		err = filtered.Remove("test.txt")

		Expect(err).NotTo(HaveOccurred())
	})

	It("should allow removing a directory", func() {
		fs := afero.NewMemMapFs()
		err := fs.Mkdir("test", os.ModeDir)
		Expect(err).NotTo(HaveOccurred())
		filtered := filter.NewFs(fs, func(name string) bool {
			return name == "not-test.txt"
		})

		err = filtered.Remove("test")

		Expect(err).NotTo(HaveOccurred())
	})

	It("should not allow stat-ing a filtered file", func() {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		filtered := filter.NewFs(fs, func(name string) bool {
			return name == "not-test.txt"
		})

		_, err = filtered.Stat("test.txt")

		Expect(err).To(MatchError(syscall.ENOENT))
	})

	It("should allow stat-ing a non-filtered file", func() {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		filtered := filter.NewFs(fs, func(name string) bool {
			return name == "test.txt"
		})

		_, err = filtered.Stat("test.txt")

		Expect(err).NotTo(HaveOccurred())
	})

	It("should walk properly", func() {
		base := afero.NewMemMapFs()
		Expect(base.Mkdir("test", os.ModePerm)).To(Succeed())
		err := afero.WriteFile(base, "test/file.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		filtered := filter.NewFs(base, func(s string) bool {
			return filepath.Ext(s) != ".txt"
		})
		paths := []string{}

		err = afero.Walk(filtered, "", func(path string, info fs.FileInfo, err error) error {
			paths = append(paths, path)
			return err
		})

		Expect(err).NotTo(HaveOccurred())
		Expect(paths).To(ConsistOf("", "test"))
	})
})
