package gitignore_test

import (
	"bytes"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	ignore "github.com/sabhiram/go-gitignore"
	"github.com/spf13/afero"
	"github.com/unmango/aferox/gitignore"
	"github.com/unmango/aferox/testing"
)

type ignoreStub string

func (s ignoreStub) MatchesPathHow(f string) (bool, *ignore.IgnorePattern) {
	return string(s) == f, nil
}

func (s ignoreStub) MatchesPath(p string) bool {
	return string(s) == p
}

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("read error")
}

var _ = Describe("Fs", func() {
	var base afero.Fs

	BeforeEach(func() {
		base = afero.NewMemMapFs()
	})

	Describe("NewFsFromLines", func() {
		It("should ignore pattern", func() {
			err := afero.WriteFile(base, "blah.txt", []byte("fdh"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			fs := gitignore.NewFsFromLines(base, "*.txt")

			_, err = fs.Stat("blah.txt")
			Expect(err).To(MatchError(os.ErrNotExist))
		})

		It("should open unignored files", func() {
			err := afero.WriteFile(base, "blah.txt", []byte("fdh"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			fs := gitignore.NewFsFromLines(base, "*.blah")

			_, err = fs.Stat("blah.txt")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("NewFsFromIgnore", func() {
		It("should ignore pattern", func() {
			err := afero.WriteFile(base, "blah.txt", []byte("fdh"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			fs := gitignore.NewFsFromIgnore(base, ignoreStub("blah.txt"))

			_, err = fs.Stat("blah.txt")
			Expect(err).To(MatchError(os.ErrNotExist))
		})

		It("should open unignored files", func() {
			err := afero.WriteFile(base, "blah.txt", []byte("fdh"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			fs := gitignore.NewFsFromIgnore(base, ignoreStub("txt.blah"))

			_, err = fs.Stat("blah.txt")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("NewFsFromReader", func() {
		It("should ignore pattern", func() {
			buf := bytes.NewBufferString("*.txt")
			err := afero.WriteFile(base, "blah.txt", []byte("fdh"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			fs, err := gitignore.NewFsFromReader(base, buf)

			Expect(err).NotTo(HaveOccurred())
			_, err = fs.Stat("blah.txt")
			Expect(err).To(MatchError(os.ErrNotExist))
		})

		It("should open unignored files", func() {
			buf := bytes.NewBufferString("*.blah")
			err := afero.WriteFile(base, "blah.txt", []byte("fdh"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			fs, err := gitignore.NewFsFromReader(base, buf)

			Expect(err).NotTo(HaveOccurred())
			_, err = fs.Stat("blah.txt")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return read errors", func() {
			buf := &errorReader{}

			fs, err := gitignore.NewFsFromReader(base, buf)

			Expect(err).To(MatchError("scanning ignore lines: read error"))
			Expect(fs).To(BeNil())
		})
	})

	Describe("NewFsFromFile", func() {
		It("should ignore pattern", func() {
			err := afero.WriteFile(base, "git.ignore", []byte("*.txt"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
			err = afero.WriteFile(base, "blah.txt", []byte("fdh"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			fs, err := gitignore.NewFsFromFile(base, "git.ignore")

			Expect(err).NotTo(HaveOccurred())
			_, err = fs.Stat("blah.txt")
			Expect(err).To(MatchError(os.ErrNotExist))
		})

		It("should open unignored files", func() {
			err := afero.WriteFile(base, "git.ignore", []byte("*.blah"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
			err = afero.WriteFile(base, "blah.txt", []byte("fdh"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			fs, err := gitignore.NewFsFromFile(base, "git.ignore")

			Expect(err).NotTo(HaveOccurred())
			_, err = fs.Stat("blah.txt")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return open errors", func() {
			base = &testing.ErrorFs{
				OpenErr: fmt.Errorf("test error"),
			}

			fs, err := gitignore.NewFsFromFile(base, "git.ignore")

			Expect(err).To(MatchError("opening ignore file: test error"))
			Expect(fs).To(BeNil())
		})
	})

	Describe("OpenDefault", func() {
		It("should ignore pattern", func() {
			err := afero.WriteFile(base, ".gitignore", []byte("*.txt"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
			err = afero.WriteFile(base, "blah.txt", []byte("fdh"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			fs, err := gitignore.OpenDefault(base)

			Expect(err).NotTo(HaveOccurred())
			_, err = fs.Stat("blah.txt")
			Expect(err).To(MatchError(os.ErrNotExist))
		})

		It("should open unignored files", func() {
			err := afero.WriteFile(base, ".gitignore", []byte("*.blah"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
			err = afero.WriteFile(base, "blah.txt", []byte("fdh"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			fs, err := gitignore.OpenDefault(base)

			Expect(err).NotTo(HaveOccurred())
			_, err = fs.Stat("blah.txt")
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
