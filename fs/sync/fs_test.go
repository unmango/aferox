package sync_test

import (
	"context"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	. "github.com/unmango/go/testing/matcher"

	"github.com/unmango/go/fs/sync"
)

var _ = Describe("Fs", func() {
	It("should read from base fs", func() {
		var (
			base  = afero.NewMemMapFs()
			layer = afero.NewMemMapFs()
		)
		err := afero.WriteFile(base, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		fs, _ := sync.NewFs(base, layer)

		Expect(fs).NotTo(BeNil())
		data, err := afero.ReadFile(fs, "test.txt")
		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).To(Equal("testing"))
	})

	It("should write to base fs", func() {
		var (
			base  = afero.NewMemMapFs()
			layer = afero.NewMemMapFs()
		)

		fs, _ := sync.NewFs(base, layer)

		Expect(fs).NotTo(BeNil())
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		Expect(base).To(ContainFileWithBytes("test.txt", []byte("testing")))
	})

	It("should not error by default", func(ctx context.Context) {
		var (
			base  = afero.NewMemMapFs()
			layer = afero.NewMemMapFs()
		)

		_, sync := sync.NewFs(base, layer)

		Expect(sync(ctx)).To(Succeed())
	})

	It("should copy layer into base", func(ctx context.Context) {
		var (
			base  = afero.NewMemMapFs()
			layer = afero.NewMemMapFs()
		)
		err := afero.WriteFile(layer, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		_, sync := sync.NewFs(base, layer)

		Expect(sync(ctx)).To(Succeed())
		Expect(base).To(ContainFileWithBytes("test.txt", []byte("testing")))
	})

	It("should use the given strategy", func(ctx context.Context) {
		var (
			base     = afero.NewMemMapFs()
			layer    = afero.NewMemMapFs()
			sentinel = false
		)
		strat := func(context.Context, afero.Fs, afero.Fs) error {
			sentinel = true
			return nil
		}

		_, sync := sync.NewFs(base, layer, sync.WithStrategy(strat))

		Expect(sync(ctx)).To(Succeed())
		Expect(sentinel).To(BeTrueBecause("the strategy is used"))
	})
})
