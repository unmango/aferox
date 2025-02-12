package aferox_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"

	"github.com/unmango/aferox"
	"github.com/unmango/go/slices"
)

var _ = Describe("Iter", func() {
	It("should iterate over an empty fs", func() {
		fs := afero.NewMemMapFs()

		seq := aferox.Iter(fs, "")

		a, b, c := slices.Collect3(seq)
		Expect(a).To(ConsistOf("")) // root dir
		Expect(b).To(HaveLen(1))
		Expect(b[0].Name()).To(Equal(""))
		Expect(c).To(ConsistOf(nil))
	})

	It("should skip root when iterating over an empty fs", func() {
		fs := afero.NewMemMapFs()

		seq := aferox.Iter(fs, "", aferox.SkipDirs)

		a, b, c := slices.Collect3(seq)
		Expect(a).To(BeEmpty())
		Expect(b).To(BeEmpty())
		Expect(c).To(BeEmpty())
	})

	It("should iterate over files", func() {
		fs := afero.NewMemMapFs()
		_, err := fs.Create("test.txt")
		Expect(err).NotTo(HaveOccurred())

		seq := aferox.Iter(fs, "")

		a, b, c := slices.Collect3(seq)
		Expect(a).To(ConsistOf("", "test.txt"))
		Expect(b).To(HaveLen(2))
		Expect(b[0].Name()).To(Equal(""))
		Expect(b[1].Name()).To(Equal("test.txt"))
		Expect(c).To(ConsistOf(nil, nil))
	})

	It("should iterate over directories", func() {
		fs := afero.NewMemMapFs()
		err := fs.Mkdir("test", os.ModeDir)
		Expect(err).NotTo(HaveOccurred())

		seq := aferox.Iter(fs, "")

		a, b, c := slices.Collect3(seq)
		Expect(a).To(ConsistOf("", "test"))
		Expect(b).To(HaveLen(2))
		Expect(b[0].Name()).To(Equal(""))
		Expect(b[1].Name()).To(Equal("test"))
		Expect(b[1].IsDir()).To(BeTrue())
		Expect(c).To(ConsistOf(nil, nil))
	})
})
