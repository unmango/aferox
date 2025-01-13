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

var _ = Describe("Adapter", func() {
	var base *testing.ContextFs

	BeforeEach(func() {
		base = &testing.ContextFs{}
	})

	It("should call base Chmod", func(ctx context.Context) {
		var (
			actualCtx   context.Context
			actualName  string
			actualMode  fs.FileMode
			expectedErr = errors.New("sentinel")
		)
		base.ChmodFunc = func(ctx context.Context, s string, fm fs.FileMode) error {
			actualCtx = ctx
			actualName = s
			actualMode = fm
			return expectedErr
		}
		fs := context.AdapterFs{base}

		err := fs.ChmodContext(ctx, "bleh", os.ModePerm)

		Expect(err).To(MatchError(expectedErr))
		Expect(actualCtx).To(BeIdenticalTo(ctx))
		Expect(actualName).To(Equal("bleh"))
		Expect(actualMode).To(Equal(os.ModePerm))
	})

	It("should call base Chown", func(ctx context.Context) {
		var (
			actualCtx   context.Context
			actualUid   int
			actualGid   int
			expectedErr = errors.New("sentinel")
		)
		base.ChownFunc = func(ctx context.Context, s string, i1, i2 int) error {
			actualCtx = ctx
			actualUid = i1
			actualGid = i2
			return expectedErr
		}
		fs := context.AdapterFs{base}

		err := fs.ChownContext(ctx, "bleh", 420, 69)

		Expect(err).To(MatchError(expectedErr))
		Expect(actualCtx).To(BeIdenticalTo(ctx))
		Expect(actualUid).To(Equal(420))
		Expect(actualGid).To(Equal(69))
	})

	It("should call base Chtimes", func(ctx context.Context) {
		var (
			actualCtx   context.Context
			actualName  string
			actualAtime time.Time
			actualMtime time.Time
			expectedErr = errors.New("sentinel")
		)
		base.ChtimesFunc = func(ctx context.Context, s string, t1, t2 time.Time) error {
			actualCtx = ctx
			actualName = s
			actualAtime = t1
			actualMtime = t2
			return expectedErr
		}
		fs := context.AdapterFs{base}

		err := fs.ChtimesContext(ctx, "bleh", time.Unix(69, 420), time.Unix(420, 69))

		Expect(err).To(MatchError(expectedErr))
		Expect(actualCtx).To(BeIdenticalTo(ctx))
		Expect(actualName).To(Equal("bleh"))
		Expect(actualAtime).To(Equal(time.Unix(69, 420)))
		Expect(actualMtime).To(Equal(time.Unix(420, 69)))
	})

	It("should call base Create", func(ctx context.Context) {
		var (
			actualCtx    context.Context
			actualName   string
			expectedFile = &testing.File{}
			expectedErr  = errors.New("sentinel")
		)
		base.CreateFunc = func(ctx context.Context, s string) (afero.File, error) {
			actualCtx = ctx
			actualName = s
			return expectedFile, expectedErr
		}
		fs := context.AdapterFs{base}

		f, err := fs.CreateContext(ctx, "bleh")

		Expect(err).To(MatchError(expectedErr))
		Expect(f).To(BeIdenticalTo(expectedFile))
		Expect(actualCtx).To(BeIdenticalTo(ctx))
		Expect(actualName).To(Equal("bleh"))
	})

	It("should call base MkdirAll", func(ctx context.Context) {
		var (
			actualCtx   context.Context
			actualName  string
			actualMode  fs.FileMode
			expectedErr = errors.New("sentinel")
		)
		base.MkdirAllFunc = func(ctx context.Context, s string, fm fs.FileMode) error {
			actualCtx = ctx
			actualName = s
			actualMode = fm
			return expectedErr
		}
		fs := context.AdapterFs{base}

		err := fs.MkdirAllContext(ctx, "bleh", os.ModeDir)

		Expect(err).To(MatchError(expectedErr))
		Expect(actualCtx).To(BeIdenticalTo(ctx))
		Expect(actualName).To(Equal("bleh"))
		Expect(actualMode).To(Equal(os.ModeDir))
	})

	It("should call base Mkdir", func(ctx context.Context) {
		var (
			actualCtx   context.Context
			actualName  string
			actualMode  fs.FileMode
			expectedErr = errors.New("sentinel")
		)
		base.MkdirFunc = func(ctx context.Context, s string, fm fs.FileMode) error {
			actualCtx = ctx
			actualName = s
			actualMode = fm
			return expectedErr
		}
		fs := context.AdapterFs{base}

		err := fs.MkdirContext(ctx, "bleh", os.ModeDir)

		Expect(err).To(MatchError(expectedErr))
		Expect(actualCtx).To(BeIdenticalTo(ctx))
		Expect(actualName).To(Equal("bleh"))
		Expect(actualMode).To(Equal(os.ModeDir))
	})

	It("should call base Open", func(ctx context.Context) {
		var (
			actualCtx    context.Context
			actualName   string
			expectedFile = &testing.File{}
			expectedErr  = errors.New("sentinel")
		)
		base.OpenFunc = func(ctx context.Context, s string) (afero.File, error) {
			actualCtx = ctx
			actualName = s
			return expectedFile, expectedErr
		}
		fs := context.AdapterFs{base}

		f, err := fs.OpenContext(ctx, "bleh")

		Expect(err).To(MatchError(expectedErr))
		Expect(f).To(BeIdenticalTo(expectedFile))
		Expect(actualCtx).To(BeIdenticalTo(ctx))
		Expect(actualName).To(Equal("bleh"))
	})

	It("should call base OpenFile", func(ctx context.Context) {
		var (
			actualCtx    context.Context
			actualName   string
			actualFlag   int
			actualMode   fs.FileMode
			expectedFile = &testing.File{}
			expectedErr  = errors.New("sentinel")
		)
		base.OpenFileFunc = func(ctx context.Context, s string, i int, fm fs.FileMode) (afero.File, error) {
			actualCtx = ctx
			actualName = s
			actualFlag = i
			actualMode = fm
			return expectedFile, expectedErr
		}
		fs := context.AdapterFs{base}

		f, err := fs.OpenFileContext(ctx, "bleh", 69, os.ModePerm)

		Expect(err).To(MatchError(expectedErr))
		Expect(f).To(BeIdenticalTo(expectedFile))
		Expect(actualCtx).To(BeIdenticalTo(ctx))
		Expect(actualName).To(Equal("bleh"))
		Expect(actualFlag).To(Equal(69))
		Expect(actualMode).To(Equal(os.ModePerm))
	})

	It("should call base RemoveAll", func(ctx context.Context) {
		var (
			actualCtx   context.Context
			actualName  string
			expectedErr = errors.New("sentinel")
		)
		base.RemoveAllFunc = func(ctx context.Context, s string) error {
			actualCtx = ctx
			actualName = s
			return expectedErr
		}
		fs := context.AdapterFs{base}

		err := fs.RemoveAllContext(ctx, "bleh")

		Expect(err).To(MatchError(expectedErr))
		Expect(actualCtx).To(BeIdenticalTo(ctx))
		Expect(actualName).To(Equal("bleh"))
	})

	It("should call base Remove", func(ctx context.Context) {
		var (
			actualCtx   context.Context
			actualName  string
			expectedErr = errors.New("sentinel")
		)
		base.RemoveFunc = func(ctx context.Context, s string) error {
			actualCtx = ctx
			actualName = s
			return expectedErr
		}
		fs := context.AdapterFs{base}

		err := fs.RemoveContext(ctx, "bleh")

		Expect(err).To(MatchError(expectedErr))
		Expect(actualCtx).To(BeIdenticalTo(ctx))
		Expect(actualName).To(Equal("bleh"))
	})

	It("should call base Rename", func(ctx context.Context) {
		var (
			actualCtx   context.Context
			actualOld   string
			actualNew   string
			expectedErr = errors.New("sentinel")
		)
		base.RenameFunc = func(ctx context.Context, s1, s2 string) error {
			actualCtx = ctx
			actualOld = s1
			actualNew = s2
			return expectedErr
		}
		fs := context.AdapterFs{base}

		err := fs.RenameContext(ctx, "bleh", "blah")

		Expect(err).To(MatchError(expectedErr))
		Expect(actualCtx).To(BeIdenticalTo(ctx))
		Expect(actualOld).To(Equal("bleh"))
		Expect(actualNew).To(Equal("blah"))
	})

	It("should call base Stat", func(ctx context.Context) {
		var (
			actualCtx    context.Context
			actualName   string
			expectedInfo = &testing.FileInfo{}
			expectedErr  = errors.New("sentinel")
		)
		base.StatFunc = func(ctx context.Context, s string) (fs.FileInfo, error) {
			actualCtx = ctx
			actualName = s
			return expectedInfo, expectedErr
		}
		fs := context.AdapterFs{base}

		i, err := fs.StatContext(ctx, "bleh")

		Expect(err).To(MatchError(expectedErr))
		Expect(actualCtx).To(BeIdenticalTo(ctx))
		Expect(actualName).To(Equal("bleh"))
		Expect(i).To(BeIdenticalTo(expectedInfo))
	})
})
