package context_test

import (
	"errors"
	"io/fs"
	"os"
	"syscall"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"

	"github.com/unmango/aferox/context"
	"github.com/unmango/aferox/testing"
)

var _ = Describe("Discard", func() {
	var base *testing.Fs

	BeforeEach(func() {
		base = &testing.Fs{}
	})

	Describe("Chmod", func() {
		It("should call base", func() {
			var (
				actualName string
				actualMode fs.FileMode
			)
			base.ChmodFunc = func(s string, fm fs.FileMode) error {
				actualName = s
				actualMode = fm
				return nil
			}
			fs := context.Discard(base)

			err := fs.Chmod("blah", os.ModePerm)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualName).To(Equal("blah"))
			Expect(actualMode).To(Equal(os.ModePerm))
		})

		It("should discard context", func(ctx context.Context) {
			var (
				actualName string
				actualMode fs.FileMode
			)
			base.ChmodFunc = func(s string, fm fs.FileMode) error {
				actualName = s
				actualMode = fm
				return nil
			}
			fs := context.Discard(base)

			err := fs.ChmodContext(ctx, "blah", os.ModePerm)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualName).To(Equal("blah"))
			Expect(actualMode).To(Equal(os.ModePerm))
		})

		It("should return the error returned by base", func() {
			expected := errors.New("sentinel")
			base.ChmodFunc = func(s string, fm fs.FileMode) error {
				return expected
			}
			fs := context.Discard(base)

			err := fs.Chmod("blah", os.ModePerm)

			Expect(err).To(MatchError(expected))
		})

		It("should return the error returned by base with context", func(ctx context.Context) {
			expected := errors.New("sentinel")
			base.ChmodFunc = func(s string, fm fs.FileMode) error {
				return expected
			}
			fs := context.Discard(base)

			err := fs.ChmodContext(ctx, "blah", os.ModePerm)

			Expect(err).To(MatchError(expected))
		})
	})

	Describe("Chown", func() {
		It("should call base", func() {
			var (
				actualName string
				actualUid  int
				actualGid  int
			)
			base.ChownFunc = func(s string, i1, i2 int) error {
				actualName = s
				actualUid = i1
				actualGid = i2
				return nil
			}
			fs := context.Discard(base)

			err := fs.Chown("blah", 69, 420)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualName).To(Equal("blah"))
			Expect(actualUid).To(Equal(69))
			Expect(actualGid).To(Equal(420))
		})

		It("should discard context", func(ctx context.Context) {
			var (
				actualName string
				actualUid  int
				actualGid  int
			)
			base.ChownFunc = func(s string, i1, i2 int) error {
				actualName = s
				actualUid = i1
				actualGid = i2
				return nil
			}
			fs := context.Discard(base)

			err := fs.ChownContext(ctx, "blah", 69, 420)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualName).To(Equal("blah"))
			Expect(actualUid).To(Equal(69))
			Expect(actualGid).To(Equal(420))
		})

		It("should return the error returned by base", func() {
			expected := errors.New("sentinel")
			base.ChownFunc = func(s string, i1, i2 int) error {
				return expected
			}
			fs := context.Discard(base)

			err := fs.Chown("blah", 69, 420)

			Expect(err).To(MatchError(expected))
		})

		It("should return the error returned by base with context", func(ctx context.Context) {
			expected := errors.New("sentinel")
			base.ChownFunc = func(s string, i1, i2 int) error {
				return expected
			}
			fs := context.Discard(base)

			err := fs.ChownContext(ctx, "blah", 69, 420)

			Expect(err).To(MatchError(expected))
		})
	})

	Describe("Chtimes", func() {
		It("should call base", func() {
			var (
				actualName  string
				actualAtime time.Time
				actualMtime time.Time
				atime       = time.Now()
				mtime       = time.Now()
			)
			base.ChtimesFunc = func(s string, t1, t2 time.Time) error {
				actualName = s
				actualAtime = t1
				actualMtime = t2
				return nil
			}
			fs := context.Discard(base)

			err := fs.Chtimes("blah", atime, mtime)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualName).To(Equal("blah"))
			Expect(actualAtime).To(Equal(atime))
			Expect(actualMtime).To(Equal(mtime))
		})

		It("should discard context", func(ctx context.Context) {
			var (
				actualName  string
				actualAtime time.Time
				actualMtime time.Time
				atime       = time.Now()
				mtime       = time.Now()
			)
			base.ChtimesFunc = func(s string, t1, t2 time.Time) error {
				actualName = s
				actualAtime = t1
				actualMtime = t2
				return nil
			}
			fs := context.Discard(base)

			err := fs.ChtimesContext(ctx, "blah", atime, mtime)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualName).To(Equal("blah"))
			Expect(actualAtime).To(Equal(atime))
			Expect(actualMtime).To(Equal(mtime))
		})

		It("should return the error returned by base", func() {
			expected := errors.New("sentinel")
			base.ChtimesFunc = func(s string, t1, t2 time.Time) error {
				return expected
			}
			fs := context.Discard(base)

			err := fs.Chtimes("blah", time.Time{}, time.Time{})

			Expect(err).To(MatchError(expected))
		})

		It("should return the error returned by base with context", func(ctx context.Context) {
			expected := errors.New("sentinel")
			base.ChtimesFunc = func(s string, t1, t2 time.Time) error {
				return expected
			}
			fs := context.Discard(base)

			err := fs.ChtimesContext(ctx, "blah", time.Time{}, time.Time{})

			Expect(err).To(MatchError(expected))
		})
	})

	Describe("Create", func() {
		It("should call base", func() {
			var (
				actualName   string
				expectedFile = &testing.File{}
			)
			base.CreateFunc = func(s string) (afero.File, error) {
				actualName = s
				return expectedFile, nil
			}
			fs := context.Discard(base)

			f, err := fs.Create("blah")

			Expect(err).NotTo(HaveOccurred())
			Expect(actualName).To(Equal("blah"))
			Expect(f).To(BeIdenticalTo(expectedFile))
		})

		It("should discard context", func(ctx context.Context) {
			var (
				actualName   string
				expectedFile = &testing.File{}
			)
			base.CreateFunc = func(s string) (afero.File, error) {
				actualName = s
				return expectedFile, nil
			}
			fs := context.Discard(base)

			f, err := fs.CreateContext(ctx, "blah")

			Expect(err).NotTo(HaveOccurred())
			Expect(actualName).To(Equal("blah"))
			Expect(f).To(BeIdenticalTo(expectedFile))
		})

		It("should return the error returned by base", func() {
			expected := errors.New("sentinel")
			base.CreateFunc = func(s string) (afero.File, error) {
				return nil, expected
			}
			fs := context.Discard(base)

			_, err := fs.Create("blah")

			Expect(err).To(MatchError(expected))
		})

		It("should return the error returned by base with context", func(ctx context.Context) {
			expected := errors.New("sentinel")
			base.CreateFunc = func(s string) (afero.File, error) {
				return nil, expected
			}
			fs := context.Discard(base)

			_, err := fs.CreateContext(ctx, "blah")

			Expect(err).To(MatchError(expected))
		})
	})

	Describe("MkdirAll", func() {
		It("should call base", func() {
			var (
				actualName string
				actualMode fs.FileMode
			)
			base.MkdirAllFunc = func(s string, fm fs.FileMode) error {
				actualName = s
				actualMode = fm
				return nil
			}
			fs := context.Discard(base)

			err := fs.MkdirAll("blah", os.ModeDir)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualName).To(Equal("blah"))
			Expect(actualMode).To(Equal(os.ModeDir))
		})

		It("should discard context", func(ctx context.Context) {
			var (
				actualName string
				actualMode fs.FileMode
			)
			base.MkdirAllFunc = func(s string, fm fs.FileMode) error {
				actualName = s
				actualMode = fm
				return nil
			}
			fs := context.Discard(base)

			err := fs.MkdirAllContext(ctx, "blah", os.ModeDir)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualName).To(Equal("blah"))
			Expect(actualMode).To(Equal(os.ModeDir))
		})

		It("should return the error returned by base", func() {
			expected := errors.New("sentinel")
			base.MkdirAllFunc = func(s string, fm fs.FileMode) error {
				return expected
			}
			fs := context.Discard(base)

			err := fs.MkdirAll("blah", os.ModeDir)

			Expect(err).To(MatchError(expected))
		})

		It("should return the error returned by base with context", func(ctx context.Context) {
			expected := errors.New("sentinel")
			base.MkdirAllFunc = func(s string, fm fs.FileMode) error {
				return expected
			}
			fs := context.Discard(base)

			err := fs.MkdirAllContext(ctx, "blah", os.ModeDir)

			Expect(err).To(MatchError(expected))
		})
	})

	Describe("Mkdir", func() {
		It("should call base", func() {
			var (
				actualName string
				actualMode fs.FileMode
			)
			base.MkdirFunc = func(s string, fm fs.FileMode) error {
				actualName = s
				actualMode = fm
				return nil
			}
			fs := context.Discard(base)

			err := fs.Mkdir("blah", os.ModeDir)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualName).To(Equal("blah"))
			Expect(actualMode).To(Equal(os.ModeDir))
		})

		It("should discard context", func(ctx context.Context) {
			var (
				actualName string
				actualMode fs.FileMode
			)
			base.MkdirFunc = func(s string, fm fs.FileMode) error {
				actualName = s
				actualMode = fm
				return nil
			}
			fs := context.Discard(base)

			err := fs.MkdirContext(ctx, "blah", os.ModeDir)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualName).To(Equal("blah"))
			Expect(actualMode).To(Equal(os.ModeDir))
		})

		It("should return the error returned by base", func() {
			expected := errors.New("sentinel")
			base.MkdirFunc = func(s string, fm fs.FileMode) error {
				return expected
			}
			fs := context.Discard(base)

			err := fs.Mkdir("blah", os.ModeDir)

			Expect(err).To(MatchError(expected))
		})

		It("should return the error returned by base with context", func(ctx context.Context) {
			expected := errors.New("sentinel")
			base.MkdirFunc = func(s string, fm fs.FileMode) error {
				return expected
			}
			fs := context.Discard(base)

			err := fs.MkdirContext(ctx, "blah", os.ModeDir)

			Expect(err).To(MatchError(expected))
		})
	})

	Describe("Open", func() {
		It("should call base", func() {
			var (
				actualName   string
				expectedFile = &testing.File{}
			)
			base.OpenFunc = func(s string) (afero.File, error) {
				actualName = s
				return expectedFile, nil
			}
			fs := context.Discard(base)

			f, err := fs.Open("blah")

			Expect(err).NotTo(HaveOccurred())
			Expect(actualName).To(Equal("blah"))
			Expect(f).To(BeIdenticalTo(expectedFile))
		})

		It("should discard context", func(ctx context.Context) {
			var (
				actualName   string
				expectedFile = &testing.File{}
			)
			base.OpenFunc = func(s string) (afero.File, error) {
				actualName = s
				return expectedFile, nil
			}
			fs := context.Discard(base)

			f, err := fs.OpenContext(ctx, "blah")

			Expect(err).NotTo(HaveOccurred())
			Expect(actualName).To(Equal("blah"))
			Expect(f).To(BeIdenticalTo(expectedFile))
		})

		It("should return the error returned by base", func() {
			expected := errors.New("sentinel")
			base.OpenFunc = func(s string) (afero.File, error) {
				return nil, expected
			}
			fs := context.Discard(base)

			_, err := fs.Open("blah")

			Expect(err).To(MatchError(expected))
		})

		It("should return the error returned by base with context", func(ctx context.Context) {
			expected := errors.New("sentinel")
			base.OpenFunc = func(s string) (afero.File, error) {
				return nil, expected
			}
			fs := context.Discard(base)

			_, err := fs.OpenContext(ctx, "blah")

			Expect(err).To(MatchError(expected))
		})
	})

	Describe("OpenFile", func() {
		It("should call base", func() {
			var (
				actualName   string
				actualFlag   int
				actualMode   fs.FileMode
				expectedFile = &testing.File{}
			)
			base.OpenFileFunc = func(s string, i int, fm fs.FileMode) (afero.File, error) {
				actualName = s
				actualFlag = i
				actualMode = fm
				return expectedFile, nil
			}
			fs := context.Discard(base)

			f, err := fs.OpenFile("blah", syscall.O_APPEND, os.ModePerm)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualName).To(Equal("blah"))
			Expect(actualFlag).To(Equal(syscall.O_APPEND))
			Expect(actualMode).To(Equal(os.ModePerm))
			Expect(f).To(BeIdenticalTo(expectedFile))
		})

		It("should discard context", func(ctx context.Context) {
			var (
				actualName   string
				actualFlag   int
				actualMode   fs.FileMode
				expectedFile = &testing.File{}
			)
			base.OpenFileFunc = func(s string, i int, fm fs.FileMode) (afero.File, error) {
				actualName = s
				actualFlag = i
				actualMode = fm
				return expectedFile, nil
			}
			fs := context.Discard(base)

			f, err := fs.OpenFileContext(ctx, "blah", syscall.O_APPEND, os.ModePerm)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualName).To(Equal("blah"))
			Expect(actualFlag).To(Equal(syscall.O_APPEND))
			Expect(actualMode).To(Equal(os.ModePerm))
			Expect(f).To(BeIdenticalTo(expectedFile))
		})

		It("should return the error returned by base", func() {
			expected := errors.New("sentinel")
			base.OpenFileFunc = func(s string, i int, fm fs.FileMode) (afero.File, error) {
				return nil, expected
			}
			fs := context.Discard(base)

			_, err := fs.OpenFile("blah", syscall.O_APPEND, os.ModePerm)

			Expect(err).To(MatchError(expected))
		})

		It("should return the error returned by base with context", func(ctx context.Context) {
			expected := errors.New("sentinel")
			base.OpenFileFunc = func(s string, i int, fm fs.FileMode) (afero.File, error) {
				return nil, expected
			}
			fs := context.Discard(base)

			_, err := fs.OpenFileContext(ctx, "blah", syscall.O_APPEND, os.ModePerm)

			Expect(err).To(MatchError(expected))
		})
	})

	Describe("Remove", func() {
		It("should call base", func() {
			var actualName string
			base.RemoveFunc = func(s string) error {
				actualName = s
				return nil
			}
			fs := context.Discard(base)

			err := fs.Remove("blah")

			Expect(err).NotTo(HaveOccurred())
			Expect(actualName).To(Equal("blah"))
		})

		It("should discard context", func(ctx context.Context) {
			var actualName string
			base.RemoveFunc = func(s string) error {
				actualName = s
				return nil
			}
			fs := context.Discard(base)

			err := fs.RemoveContext(ctx, "blah")

			Expect(err).NotTo(HaveOccurred())
			Expect(actualName).To(Equal("blah"))
		})

		It("should return the error returned by base", func() {
			expected := errors.New("sentinel")
			base.RemoveFunc = func(s string) error {
				return expected
			}
			fs := context.Discard(base)

			err := fs.Remove("blah")

			Expect(err).To(MatchError(expected))
		})

		It("should return the error returned by base with context", func(ctx context.Context) {
			expected := errors.New("sentinel")
			base.RemoveFunc = func(s string) error {
				return expected
			}
			fs := context.Discard(base)

			err := fs.RemoveContext(ctx, "blah")

			Expect(err).To(MatchError(expected))
		})
	})

	Describe("RemoveAll", func() {
		It("should call base", func() {
			var actualName string
			base.RemoveAllFunc = func(s string) error {
				actualName = s
				return nil
			}
			fs := context.Discard(base)

			err := fs.RemoveAll("blah")

			Expect(err).NotTo(HaveOccurred())
			Expect(actualName).To(Equal("blah"))
		})

		It("should discard context", func(ctx context.Context) {
			var actualName string
			base.RemoveAllFunc = func(s string) error {
				actualName = s
				return nil
			}
			fs := context.Discard(base)

			err := fs.RemoveAllContext(ctx, "blah")

			Expect(err).NotTo(HaveOccurred())
			Expect(actualName).To(Equal("blah"))
		})

		It("should return the error returned by base", func() {
			expected := errors.New("sentinel")
			base.RemoveAllFunc = func(s string) error {
				return expected
			}
			fs := context.Discard(base)

			err := fs.RemoveAll("blah")

			Expect(err).To(MatchError(expected))
		})

		It("should return the error returned by base with context", func(ctx context.Context) {
			expected := errors.New("sentinel")
			base.RemoveAllFunc = func(s string) error {
				return expected
			}
			fs := context.Discard(base)

			err := fs.RemoveAllContext(ctx, "blah")

			Expect(err).To(MatchError(expected))
		})
	})

	Describe("Rename", func() {
		It("should call base", func() {
			var (
				actualOld string
				actualNew string
			)
			base.RenameFunc = func(s1, s2 string) error {
				actualOld = s1
				actualNew = s2
				return nil
			}
			fs := context.Discard(base)

			err := fs.Rename("blah", "bleh")

			Expect(err).NotTo(HaveOccurred())
			Expect(actualOld).To(Equal("blah"))
			Expect(actualNew).To(Equal("bleh"))
		})

		It("should discard context", func(ctx context.Context) {
			var (
				actualOld string
				actualNew string
			)
			base.RenameFunc = func(s1, s2 string) error {
				actualOld = s1
				actualNew = s2
				return nil
			}
			fs := context.Discard(base)

			err := fs.RenameContext(ctx, "blah", "bleh")

			Expect(err).NotTo(HaveOccurred())
			Expect(actualOld).To(Equal("blah"))
			Expect(actualNew).To(Equal("bleh"))
		})

		It("should return the error returned by base", func() {
			expected := errors.New("sentinel")
			base.RenameFunc = func(s1, s2 string) error {
				return expected
			}
			fs := context.Discard(base)

			err := fs.Rename("blah", "bleh")

			Expect(err).To(MatchError(expected))
		})

		It("should return the error returned by base with context", func(ctx context.Context) {
			expected := errors.New("sentinel")
			base.RenameFunc = func(s1, s2 string) error {
				return expected
			}
			fs := context.Discard(base)

			err := fs.RenameContext(ctx, "blah", "bleh")

			Expect(err).To(MatchError(expected))
		})
	})

	Describe("Stat", func() {
		It("should call base", func() {
			var (
				actualName   string
				expectedInfo = &testing.FileInfo{}
			)
			base.StatFunc = func(s string) (fs.FileInfo, error) {
				actualName = s
				return expectedInfo, nil
			}
			fs := context.Discard(base)

			i, err := fs.Stat("blah")

			Expect(err).NotTo(HaveOccurred())
			Expect(actualName).To(Equal("blah"))
			Expect(i).To(BeIdenticalTo(expectedInfo))
		})

		It("should discard context", func(ctx context.Context) {
			var (
				actualName   string
				expectedInfo = &testing.FileInfo{}
			)
			base.StatFunc = func(s string) (fs.FileInfo, error) {
				actualName = s
				return expectedInfo, nil
			}
			fs := context.Discard(base)

			i, err := fs.StatContext(ctx, "blah")

			Expect(err).NotTo(HaveOccurred())
			Expect(actualName).To(Equal("blah"))
			Expect(i).To(BeIdenticalTo(expectedInfo))
		})

		It("should return the error returned by base", func() {
			expected := errors.New("sentinel")
			base.StatFunc = func(s string) (fs.FileInfo, error) {
				return nil, expected
			}
			fs := context.Discard(base)

			_, err := fs.Stat("blah")

			Expect(err).To(MatchError(expected))
		})

		It("should return the error returned by base with context", func(ctx context.Context) {
			expected := errors.New("sentinel")
			base.StatFunc = func(s string) (fs.FileInfo, error) {
				return nil, expected
			}
			fs := context.Discard(base)

			_, err := fs.StatContext(ctx, "blah")

			Expect(err).To(MatchError(expected))
		})
	})
})
