package ignore_test

import (
	"io/fs"
	"os"
	"syscall"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/spf13/afero"
	"github.com/unmango/aferox/ignore"
)

type ignoreStub string

func (s ignoreStub) MatchesPath(p string) bool {
	return string(s) == p
}

var _ = Describe("Fs", func() {
	It("should not allow chmod-ing a ignored file", func() {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		ignored := ignore.NewFs(fs, ignoreStub("not-test.txt"))

		err = ignored.Chmod("test.txt", os.ModeAppend)

		Expect(err).To(MatchError(syscall.ENOENT))
	})

	It("should allow chmod-ing a non-ignored file", func() {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		ignored := ignore.NewFs(fs, ignoreStub("test.txt"))

		err = ignored.Chmod("test.txt", os.ModeAppend)

		Expect(err).NotTo(HaveOccurred())
	})

	It("should not allow chown-ing a ignored file", func() {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		ignored := ignore.NewFs(fs, ignoreStub("not-test.txt"))

		err = ignored.Chown("test.txt", 1001, 1001)

		Expect(err).To(MatchError(syscall.ENOENT))
	})

	It("should allow chown-ing a non-ignored file", func() {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		ignored := ignore.NewFs(fs, ignoreStub("test.txt"))

		err = ignored.Chown("test.txt", 1001, 1001)

		Expect(err).NotTo(HaveOccurred())
	})

	It("should not allow chtimes-ing a ignored file", func() {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		ignored := ignore.NewFs(fs, ignoreStub("not-test.txt"))

		err = ignored.Chtimes("test.txt", time.Now(), time.Now())

		Expect(err).To(MatchError(syscall.ENOENT))
	})

	It("should allow chtimes-ing a non-ignored file", func() {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		ignored := ignore.NewFs(fs, ignoreStub("test.txt"))

		err = ignored.Chtimes("test.txt", time.Now(), time.Now())

		Expect(err).NotTo(HaveOccurred())
	})

	It("should not allow creating a ignored file", func() {
		fs := afero.NewMemMapFs()
		ignored := ignore.NewFs(fs, ignoreStub("not-test.txt"))

		_, err := ignored.Create("test.txt")

		Expect(err).To(MatchError(syscall.ENOENT))
	})

	It("should allow creating a non-ignored file", func() {
		fs := afero.NewMemMapFs()
		ignored := ignore.NewFs(fs, ignoreStub("test.txt"))

		_, err := ignored.Create("test.txt")

		Expect(err).NotTo(HaveOccurred())
	})

	It("should include the source filesystem name", func() {
		fs := ignore.NewFs(afero.NewMemMapFs(), nil)

		Expect(fs.Name()).To(Equal("ignore: MemMapFS"))
	})

	It("should not allow opening a ignored file", func() {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		ignored := ignore.NewFs(fs, ignoreStub("not-test.txt"))

		_, err = ignored.Open("test.txt")

		Expect(err).To(MatchError(syscall.ENOENT))
	})

	It("should allow opening a non-ignored file", func() {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		ignored := ignore.NewFs(fs, ignoreStub("test.txt"))

		_, err = ignored.Open("test.txt")

		Expect(err).NotTo(HaveOccurred())
	})

	It("should not allow open-file-ing a ignored file", func() {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		ignored := ignore.NewFs(fs, ignoreStub("not-test.txt"))

		_, err = ignored.OpenFile("test.txt", 69, os.ModeAppend)

		Expect(err).To(MatchError(syscall.ENOENT))
	})

	It("should allow open-file-ing a non-ignored file", func() {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		ignored := ignore.NewFs(fs, ignoreStub("test.txt"))

		_, err = ignored.OpenFile("test.txt", 69, os.ModeAppend)

		Expect(err).NotTo(HaveOccurred())
	})

	It("should not allow removing a ignored file", func() {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		ignored := ignore.NewFs(fs, ignoreStub("not-test.txt"))

		err = ignored.Remove("test.txt")

		Expect(err).To(MatchError(syscall.ENOENT))
	})

	It("should allow removing a non-ignored file", func() {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		ignored := ignore.NewFs(fs, ignoreStub("test.txt"))

		err = ignored.Remove("test.txt")

		Expect(err).NotTo(HaveOccurred())
	})

	It("should allow removing a directory", func() {
		fs := afero.NewMemMapFs()
		err := fs.Mkdir("test", os.ModeDir)
		Expect(err).NotTo(HaveOccurred())
		ignored := ignore.NewFs(fs, ignoreStub("not-test.txt"))

		err = ignored.Remove("test")

		Expect(err).NotTo(HaveOccurred())
	})

	It("should not allow stat-ing a ignored file", func() {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		ignored := ignore.NewFs(fs, ignoreStub("not-test.txt"))

		_, err = ignored.Stat("test.txt")

		Expect(err).To(MatchError(syscall.ENOENT))
	})

	It("should allow stat-ing a non-ignored file", func() {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		ignored := ignore.NewFs(fs, ignoreStub("test.txt"))

		_, err = ignored.Stat("test.txt")

		Expect(err).NotTo(HaveOccurred())
	})

	It("should walk properly", func() {
		base := afero.NewMemMapFs()
		Expect(base.Mkdir("test", os.ModePerm)).To(Succeed())
		err := afero.WriteFile(base, "test/file.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		ignored := ignore.NewFs(base, ignoreStub("test/file.txt"))
		paths := []string{}

		err = afero.Walk(ignored, "", func(path string, info fs.FileInfo, err error) error {
			paths = append(paths, path)
			return err
		})

		Expect(err).NotTo(HaveOccurred())
		Expect(paths).To(ConsistOf("", "test"))
	})
})
