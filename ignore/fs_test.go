package ignore_test

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
	"github.com/unmango/aferox/ignore"
	"github.com/unmango/aferox/op"
)

type ignoreStub string

func (s ignoreStub) MatchesPath(p string) bool {
	return string(s) == p
}

var _ = Describe("Fs", func() {
	When("a single file exists", func() {
		var baseFs afero.Fs

		BeforeEach(func() {
			baseFs = afero.NewMemMapFs()
			err := afero.WriteFile(baseFs, "test.txt", []byte("testing"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
		})

		When("the file is ignored", func() {
			var ignoreFs afero.Fs

			BeforeEach(func() {
				ignoreFs = ignore.NewFs(baseFs, ignoreStub("test.txt"))
			})

			It("should not allow chmod-ing a filtered file", func() {
				err := ignoreFs.Chmod("test.txt", os.ModeAppend)

				Expect(err).To(MatchError(syscall.ENOENT))
			})

			It("should not allow chown-ing a filtered file", func() {
				err := ignoreFs.Chown("test.txt", 1001, 1001)

				Expect(err).To(MatchError(syscall.ENOENT))
			})

			It("should not allow chtimes-ing a filtered file", func() {
				err := ignoreFs.Chtimes("test.txt", time.Now(), time.Now())

				Expect(err).To(MatchError(syscall.ENOENT))
			})

			It("should not allow creating a filtered file", func() {
				_, err := ignoreFs.Create("test.txt")

				Expect(err).To(MatchError(syscall.ENOENT))
			})

			It("should not allow opening a filtered file", func() {
				_, err := ignoreFs.Open("test.txt")

				Expect(err).To(MatchError(syscall.ENOENT))
			})

			It("should not allow open-file-ing a filtered file", func() {
				_, err := ignoreFs.OpenFile("test.txt", 69, os.ModeAppend)

				Expect(err).To(MatchError(syscall.ENOENT))
			})

			It("should not allow removing a filtered file", func() {
				err := ignoreFs.Remove("test.txt")

				Expect(err).To(MatchError(syscall.ENOENT))
			})

			It("should not allow stat-ing a filtered file", func() {
				_, err := ignoreFs.Stat("test.txt")

				Expect(err).To(MatchError(syscall.ENOENT))
			})
		})

		When("the file is not ignored", func() {
			var filtered afero.Fs

			BeforeEach(func() {
				filtered = ignore.NewFs(baseFs, ignoreStub("not-test.txt"))
			})

			It("should allow chmod-ing a non-filtered file", func() {
				err := filtered.Chmod("test.txt", os.ModeAppend)

				Expect(err).NotTo(HaveOccurred())
			})

			It("should allow chown-ing a non-filtered file", func() {
				err := filtered.Chown("test.txt", 1001, 1001)

				Expect(err).NotTo(HaveOccurred())
			})

			It("should allow chtimes-ing a non-filtered file", func() {
				err := filtered.Chtimes("test.txt", time.Now(), time.Now())

				Expect(err).NotTo(HaveOccurred())
			})

			It("should allow creating a non-filtered file", func() {
				_, err := filtered.Create("test.txt")

				Expect(err).NotTo(HaveOccurred())
			})

			It("should include the source filesystem name", func() {
				Expect(filtered.Name()).To(Equal("Filter: MemMapFS"))
			})

			It("should allow opening a non-filtered file", func() {
				_, err := filtered.Open("test.txt")

				Expect(err).NotTo(HaveOccurred())
			})

			It("should allow open-file-ing a non-filtered file", func() {
				_, err := filtered.OpenFile("test.txt", 69, os.ModeAppend)

				Expect(err).NotTo(HaveOccurred())
			})

			It("should allow removing a non-filtered file", func() {
				err := filtered.Remove("test.txt")

				Expect(err).NotTo(HaveOccurred())
			})

			It("should allow removing a directory", func() {
				err := baseFs.Mkdir("test", os.ModeDir)
				Expect(err).NotTo(HaveOccurred())

				err = filtered.Remove("test")

				Expect(err).NotTo(HaveOccurred())
			})

			It("should allow stat-ing a non-filtered file", func() {
				_, err := filtered.Stat("test.txt")

				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	It("should walk properly", func() {
		base := afero.NewMemMapFs()
		Expect(base.Mkdir("test", os.ModePerm)).To(Succeed())
		err := afero.WriteFile(base, "test/file.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		filtered := filter.NewFs(base, func(op op.Operation) bool {
			return filepath.Ext(op.Path()) != ".txt"
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
