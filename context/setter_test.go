package context_test

import (
	"errors"
	"io/fs"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"

	"github.com/unmango/aferox/context"
	"github.com/unmango/aferox/testing"
)

type noopSetter struct{}

func (noopSetter) SetContext(context.Context) {}

type setter struct {
	Ctx context.Context
}

func (c *setter) SetContext(ctx context.Context) {
	c.Ctx = ctx
}

var _ = Describe("Setter", func() {
	var base *testing.Fs

	BeforeEach(func() {
		base = &testing.Fs{}
	})

	It("should call base Chmod", func(ctx context.Context) {
		var (
			actualName  string
			actualMode  fs.FileMode
			expectedErr = errors.New("sentinel")
		)
		base.ChmodFunc = func(s string, fm fs.FileMode) error {
			actualName = s
			actualMode = fm
			return expectedErr
		}
		fs := context.WithSetterFs{noopSetter{}, base}

		err := fs.Chmod(ctx, "bleh", os.ModePerm)

		Expect(err).To(MatchError(expectedErr))
		Expect(actualName).To(Equal("bleh"))
		Expect(actualMode).To(Equal(os.ModePerm))
	})

	It("should set the context when calling Chmod", func(ctx context.Context) {
		base.ChmodFunc = func(string, fs.FileMode) error { return nil }
		s := &setter{}
		fs := context.WithSetterFs{s, base}

		err := fs.Chmod(ctx, "bleh", os.ModePerm)

		Expect(err).NotTo(HaveOccurred())
		Expect(s.Ctx).To(BeIdenticalTo(ctx))
	})

	It("should call base Chown", func(ctx context.Context) {
		var (
			actualUid   int
			actualGid   int
			expectedErr = errors.New("sentinel")
		)
		base.ChownFunc = func(s string, i1, i2 int) error {
			actualUid = i1
			actualGid = i2
			return expectedErr
		}
		fs := context.WithSetterFs{noopSetter{}, base}

		err := fs.Chown(ctx, "bleh", 420, 69)

		Expect(err).To(MatchError(expectedErr))
		Expect(actualUid).To(Equal(420))
		Expect(actualGid).To(Equal(69))
	})

	It("should set the context when calling Chown", func(ctx context.Context) {
		base.ChownFunc = func(s string, i1, i2 int) error { return nil }
		s := &setter{}
		fs := context.WithSetterFs{s, base}

		err := fs.Chown(ctx, "bleh", 69, 420)

		Expect(err).NotTo(HaveOccurred())
		Expect(s.Ctx).To(BeIdenticalTo(ctx))
	})

	It("should call base Chtimes", func(ctx context.Context) {
		var (
			actualName  string
			actualAtime time.Time
			actualMtime time.Time
			expectedErr = errors.New("sentinel")
		)
		base.ChtimesFunc = func(s string, t1, t2 time.Time) error {
			actualName = s
			actualAtime = t1
			actualMtime = t2
			return expectedErr
		}
		fs := context.WithSetterFs{noopSetter{}, base}

		err := fs.Chtimes(ctx, "bleh", time.Unix(69, 420), time.Unix(420, 69))

		Expect(err).To(MatchError(expectedErr))
		Expect(actualName).To(Equal("bleh"))
		Expect(actualAtime).To(Equal(time.Unix(69, 420)))
		Expect(actualMtime).To(Equal(time.Unix(420, 69)))
	})

	It("should set the context when calling Chtimes", func(ctx context.Context) {
		base.ChtimesFunc = func(s string, t1, t2 time.Time) error { return nil }
		s := &setter{}
		fs := context.WithSetterFs{s, base}

		err := fs.Chtimes(ctx, "bleh", time.Time{}, time.Time{})

		Expect(err).NotTo(HaveOccurred())
		Expect(s.Ctx).To(BeIdenticalTo(ctx))
	})

	It("should call base Create", func(ctx context.Context) {
		var (
			actualName   string
			expectedFile = &testing.File{}
			expectedErr  = errors.New("sentinel")
		)
		base.CreateFunc = func(s string) (afero.File, error) {
			actualName = s
			return expectedFile, expectedErr
		}
		fs := context.WithSetterFs{noopSetter{}, base}

		f, err := fs.Create(ctx, "bleh")

		Expect(err).To(MatchError(expectedErr))
		Expect(f).To(BeIdenticalTo(expectedFile))
		Expect(actualName).To(Equal("bleh"))
	})

	It("should set the context when calling Create", func(ctx context.Context) {
		base.CreateFunc = func(string) (afero.File, error) { return nil, nil }
		s := &setter{}
		fs := context.WithSetterFs{s, base}

		_, err := fs.Create(ctx, "bleh")

		Expect(err).NotTo(HaveOccurred())
		Expect(s.Ctx).To(BeIdenticalTo(ctx))
	})

	It("should call base MkdirAll", func(ctx context.Context) {
		var (
			actualName  string
			actualMode  fs.FileMode
			expectedErr = errors.New("sentinel")
		)
		base.MkdirAllFunc = func(s string, fm fs.FileMode) error {
			actualName = s
			actualMode = fm
			return expectedErr
		}
		fs := context.WithSetterFs{noopSetter{}, base}

		err := fs.MkdirAll(ctx, "bleh", os.ModeDir)

		Expect(err).To(MatchError(expectedErr))
		Expect(actualName).To(Equal("bleh"))
		Expect(actualMode).To(Equal(os.ModeDir))
	})

	It("should set the context when calling MkdirAll", func(ctx context.Context) {
		base.MkdirAllFunc = func(string, fs.FileMode) error { return nil }
		s := &setter{}
		fs := context.WithSetterFs{s, base}

		err := fs.MkdirAll(ctx, "bleh", 0)

		Expect(err).NotTo(HaveOccurred())
		Expect(s.Ctx).To(BeIdenticalTo(ctx))
	})

	It("should call base Mkdir", func(ctx context.Context) {
		var (
			actualName  string
			actualMode  fs.FileMode
			expectedErr = errors.New("sentinel")
		)
		base.MkdirFunc = func(s string, fm fs.FileMode) error {
			actualName = s
			actualMode = fm
			return expectedErr
		}
		fs := context.WithSetterFs{noopSetter{}, base}

		err := fs.Mkdir(ctx, "bleh", os.ModeDir)

		Expect(err).To(MatchError(expectedErr))
		Expect(actualName).To(Equal("bleh"))
		Expect(actualMode).To(Equal(os.ModeDir))
	})

	It("should set the context when calling Mkdir", func(ctx context.Context) {
		base.MkdirFunc = func(string, fs.FileMode) error { return nil }
		s := &setter{}
		fs := context.WithSetterFs{s, base}

		err := fs.Mkdir(ctx, "bleh", 0)

		Expect(err).NotTo(HaveOccurred())
		Expect(s.Ctx).To(BeIdenticalTo(ctx))
	})

	It("should call base Open", func(ctx context.Context) {
		var (
			actualName   string
			expectedFile = &testing.File{}
			expectedErr  = errors.New("sentinel")
		)
		base.OpenFunc = func(s string) (afero.File, error) {
			actualName = s
			return expectedFile, expectedErr
		}
		fs := context.WithSetterFs{noopSetter{}, base}

		f, err := fs.Open(ctx, "bleh")

		Expect(err).To(MatchError(expectedErr))
		Expect(f).To(BeIdenticalTo(expectedFile))
		Expect(actualName).To(Equal("bleh"))
	})

	It("should set the context when calling Open", func(ctx context.Context) {
		base.OpenFunc = func(string) (afero.File, error) { return nil, nil }
		s := &setter{}
		fs := context.WithSetterFs{s, base}

		_, err := fs.Open(ctx, "bleh")

		Expect(err).NotTo(HaveOccurred())
		Expect(s.Ctx).To(BeIdenticalTo(ctx))
	})

	It("should call base OpenFile", func(ctx context.Context) {
		var (
			actualName   string
			actualFlag   int
			actualMode   fs.FileMode
			expectedFile = &testing.File{}
			expectedErr  = errors.New("sentinel")
		)
		base.OpenFileFunc = func(s string, i int, fm fs.FileMode) (afero.File, error) {
			actualName = s
			actualFlag = i
			actualMode = fm
			return expectedFile, expectedErr
		}
		fs := context.WithSetterFs{noopSetter{}, base}

		f, err := fs.OpenFile(ctx, "bleh", 69, os.ModePerm)

		Expect(err).To(MatchError(expectedErr))
		Expect(f).To(BeIdenticalTo(expectedFile))
		Expect(actualName).To(Equal("bleh"))
		Expect(actualFlag).To(Equal(69))
		Expect(actualMode).To(Equal(os.ModePerm))
	})

	It("should set the context when calling OpenFile", func(ctx context.Context) {
		base.OpenFileFunc = func(string, int, fs.FileMode) (afero.File, error) { return nil, nil }
		s := &setter{}
		fs := context.WithSetterFs{s, base}

		_, err := fs.OpenFile(ctx, "bleh", 0, 0)

		Expect(err).NotTo(HaveOccurred())
		Expect(s.Ctx).To(BeIdenticalTo(ctx))
	})

	It("should call base RemoveAll", func(ctx context.Context) {
		var (
			actualName  string
			expectedErr = errors.New("sentinel")
		)
		base.RemoveAllFunc = func(s string) error {
			actualName = s
			return expectedErr
		}
		fs := context.WithSetterFs{noopSetter{}, base}

		err := fs.RemoveAll(ctx, "bleh")

		Expect(err).To(MatchError(expectedErr))
		Expect(actualName).To(Equal("bleh"))
	})

	It("should set the context when calling RemoveAll", func(ctx context.Context) {
		base.RemoveAllFunc = func(string) error { return nil }
		s := &setter{}
		fs := context.WithSetterFs{s, base}

		err := fs.RemoveAll(ctx, "bleh")

		Expect(err).NotTo(HaveOccurred())
		Expect(s.Ctx).To(BeIdenticalTo(ctx))
	})

	It("should call base Remove", func(ctx context.Context) {
		var (
			actualName  string
			expectedErr = errors.New("sentinel")
		)
		base.RemoveFunc = func(s string) error {
			actualName = s
			return expectedErr
		}
		fs := context.WithSetterFs{noopSetter{}, base}

		err := fs.Remove(ctx, "bleh")

		Expect(err).To(MatchError(expectedErr))
		Expect(actualName).To(Equal("bleh"))
	})

	It("should set the context when calling Remove", func(ctx context.Context) {
		base.RemoveFunc = func(string) error { return nil }
		s := &setter{}
		fs := context.WithSetterFs{s, base}

		err := fs.Remove(ctx, "bleh")

		Expect(err).NotTo(HaveOccurred())
		Expect(s.Ctx).To(BeIdenticalTo(ctx))
	})

	It("should call base Rename", func(ctx context.Context) {
		var (
			actualOld   string
			actualNew   string
			expectedErr = errors.New("sentinel")
		)
		base.RenameFunc = func(s1, s2 string) error {
			actualOld = s1
			actualNew = s2
			return expectedErr
		}
		fs := context.WithSetterFs{noopSetter{}, base}

		err := fs.Rename(ctx, "bleh", "blah")

		Expect(err).To(MatchError(expectedErr))
		Expect(actualOld).To(Equal("bleh"))
		Expect(actualNew).To(Equal("blah"))
	})

	It("should set the context when calling Rename", func(ctx context.Context) {
		base.RenameFunc = func(s1, s2 string) error { return nil }
		s := &setter{}
		fs := context.WithSetterFs{s, base}

		err := fs.Rename(ctx, "bleh", "blah")

		Expect(err).NotTo(HaveOccurred())
		Expect(s.Ctx).To(BeIdenticalTo(ctx))
	})

	It("should call base Stat", func(ctx context.Context) {
		var (
			actualName   string
			expectedInfo = &testing.FileInfo{}
			expectedErr  = errors.New("sentinel")
		)
		base.StatFunc = func(s string) (fs.FileInfo, error) {
			actualName = s
			return expectedInfo, expectedErr
		}
		fs := context.WithSetterFs{noopSetter{}, base}

		i, err := fs.Stat(ctx, "bleh")

		Expect(err).To(MatchError(expectedErr))
		Expect(actualName).To(Equal("bleh"))
		Expect(i).To(BeIdenticalTo(expectedInfo))
	})

	It("should set the context when calling Stat", func(ctx context.Context) {
		base.StatFunc = func(string) (fs.FileInfo, error) { return nil, nil }
		s := &setter{}
		fs := context.WithSetterFs{s, base}

		_, err := fs.Stat(ctx, "bleh")

		Expect(err).NotTo(HaveOccurred())
		Expect(s.Ctx).To(BeIdenticalTo(ctx))
	})
})
