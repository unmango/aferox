package tar_test

import (
	"bytes"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	aferoxtar "github.com/unmango/aferox/writer/tar"
)

var _ = Describe("FileInfo", func() {
	It("should return false for IsDir", func() {
		buf := &bytes.Buffer{}
		fs := aferoxtar.NewFs(aferoxtar.NewWriter(buf))

		info, err := fs.Stat("test.txt")
		Expect(err).ToNot(HaveOccurred())
		Expect(info.IsDir()).To(BeFalse())
	})

	It("should return zero time for ModTime", func() {
		buf := &bytes.Buffer{}
		fs := aferoxtar.NewFs(aferoxtar.NewWriter(buf))

		info, err := fs.Stat("test.txt")
		Expect(err).ToNot(HaveOccurred())
		Expect(info.ModTime().IsZero()).To(BeTrue())
	})

	It("should return ModePerm for Mode", func() {
		buf := &bytes.Buffer{}
		fs := aferoxtar.NewFs(aferoxtar.NewWriter(buf))

		info, err := fs.Stat("test.txt")
		Expect(err).ToNot(HaveOccurred())
		Expect(info.Mode()).To(Equal(os.ModePerm))
	})

	It("should return correct name", func() {
		buf := &bytes.Buffer{}
		fs := aferoxtar.NewFs(aferoxtar.NewWriter(buf))

		info, err := fs.Stat("myfile.txt")
		Expect(err).ToNot(HaveOccurred())
		Expect(info.Name()).To(Equal("myfile.txt"))
	})

	It("should return 0 for Size", func() {
		buf := &bytes.Buffer{}
		fs := aferoxtar.NewFs(aferoxtar.NewWriter(buf))

		info, err := fs.Stat("test.txt")
		Expect(err).ToNot(HaveOccurred())
		Expect(info.Size()).To(Equal(int64(0)))
	})

	It("should return nil for Sys", func() {
		buf := &bytes.Buffer{}
		fs := aferoxtar.NewFs(aferoxtar.NewWriter(buf))

		info, err := fs.Stat("test.txt")
		Expect(err).ToNot(HaveOccurred())
		Expect(info.Sys()).To(BeNil())
	})
})
