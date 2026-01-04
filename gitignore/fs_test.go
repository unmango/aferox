package ignore_test

import (
	"bytes"
	"os"

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
	var base afero.Fs

	BeforeEach(func() {
		base = afero.NewMemMapFs()
	})

	Describe("NewFsFromGitIgnoreLines", func() {
		It("should ignore pattern", func() {
			err := afero.WriteFile(base, "blah.txt", []byte("fdh"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			fs := ignore.NewFsFromGitIgnoreLines(base, "*.txt")

			_, err = fs.Stat("blah.txt")
			Expect(err).To(MatchError(os.ErrNotExist))
		})

		It("should open unignored files", func() {
			err := afero.WriteFile(base, "blah.txt", []byte("fdh"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			fs := ignore.NewFsFromGitIgnoreLines(base, "*.blah")

			_, err = fs.Stat("blah.txt")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("NewFsFromIgnore", func() {
		It("should ignore pattern", func() {
			err := afero.WriteFile(base, "blah.txt", []byte("fdh"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			fs := ignore.NewFsFromIgnore(base, ignoreStub("blah.txt"))

			_, err = fs.Stat("blah.txt")
			Expect(err).To(MatchError(os.ErrNotExist))
		})

		It("should open unignored files", func() {
			err := afero.WriteFile(base, "blah.txt", []byte("fdh"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			fs := ignore.NewFsFromIgnore(base, ignoreStub("txt.blah"))

			_, err = fs.Stat("blah.txt")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("NewFsFromGitIgnoreReader", func() {
		It("should ignore pattern", func() {
			buf := bytes.NewBufferString("*.txt")
			err := afero.WriteFile(base, "blah.txt", []byte("fdh"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			fs, err := ignore.NewFsFromGitIgnoreReader(base, buf)

			Expect(err).NotTo(HaveOccurred())
			_, err = fs.Stat("blah.txt")
			Expect(err).To(MatchError(os.ErrNotExist))
		})

		It("should open unignored files", func() {
			buf := bytes.NewBufferString("*.blah")
			err := afero.WriteFile(base, "blah.txt", []byte("fdh"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			fs, err := ignore.NewFsFromGitIgnoreReader(base, buf)

			Expect(err).NotTo(HaveOccurred())
			_, err = fs.Stat("blah.txt")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("NewFsFromGitIgnoreFile", func() {
		It("should ignore pattern", func() {
			err := afero.WriteFile(base, "git.ignore", []byte("*.txt"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
			err = afero.WriteFile(base, "blah.txt", []byte("fdh"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			fs, err := ignore.NewFsFromGitIgnoreFile(base, "git.ignore")

			Expect(err).NotTo(HaveOccurred())
			_, err = fs.Stat("blah.txt")
			Expect(err).To(MatchError(os.ErrNotExist))
		})

		It("should open unignored files", func() {
			err := afero.WriteFile(base, "git.ignore", []byte("*.blah"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
			err = afero.WriteFile(base, "blah.txt", []byte("fdh"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			fs, err := ignore.NewFsFromGitIgnoreFile(base, "git.ignore")

			Expect(err).NotTo(HaveOccurred())
			_, err = fs.Stat("blah.txt")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("OpenDefaultGitIgnore", func() {
		It("should ignore pattern", func() {
			err := afero.WriteFile(base, ".gitignore", []byte("*.txt"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
			err = afero.WriteFile(base, "blah.txt", []byte("fdh"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			fs, err := ignore.OpenDefaultGitIgnore(base)

			Expect(err).NotTo(HaveOccurred())
			_, err = fs.Stat("blah.txt")
			Expect(err).To(MatchError(os.ErrNotExist))
		})

		It("should open unignored files", func() {
			err := afero.WriteFile(base, ".gitignore", []byte("*.blah"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
			err = afero.WriteFile(base, "blah.txt", []byte("fdh"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			fs, err := ignore.OpenDefaultGitIgnore(base)

			Expect(err).NotTo(HaveOccurred())
			_, err = fs.Stat("blah.txt")
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
