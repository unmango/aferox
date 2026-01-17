package tar_test

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/unmango/aferox/writer/tar"
)

var _ = Describe("Fs", func() {
	It("should create a file", func() {
		buf := &bytes.Buffer{}
		fs := tar.NewFs(tar.NewWriter(buf))

		file, err := fs.Create("test.txt")
		Expect(err).ToNot(HaveOccurred())
		defer file.Close()

		_, err = file.WriteString("Hello, World!")
		Expect(err).ToNot(HaveOccurred())
		Expect(file.Close()).To(Succeed())

		r := tar.NewReader(buf)
		header, err := r.Next()
		Expect(err).ToNot(HaveOccurred())
		Expect(header.Name).To(Equal("test.txt"))
		Expect(header.Size).To(Equal(int64(13)))
	})

	It("should write multiple files", func() {
		buf := &bytes.Buffer{}
		fs := tar.NewFs(tar.NewWriter(buf))

		file1, err := fs.Create("file1.txt")
		Expect(err).ToNot(HaveOccurred())
		Expect(file1.Close()).To(Succeed())

		file2, err := fs.Create("file2.txt")
		Expect(err).ToNot(HaveOccurred())
		Expect(file2.Close()).To(Succeed())

		file3, err := fs.Create("file3.txt")
		Expect(err).ToNot(HaveOccurred())
		Expect(file3.Close()).To(Succeed())

		r := tar.NewReader(buf)
		header, err := r.Next()
		Expect(err).ToNot(HaveOccurred())
		Expect(header.Name).To(Equal("file1.txt"))
		Expect(header.Size).To(Equal(int64(0)))

		header, err = r.Next()
		Expect(err).ToNot(HaveOccurred())
		Expect(header.Name).To(Equal("file2.txt"))
		Expect(header.Size).To(Equal(int64(0)))

		header, err = r.Next()
		Expect(err).ToNot(HaveOccurred())
		Expect(header.Name).To(Equal("file3.txt"))
		Expect(header.Size).To(Equal(int64(0)))
	})

	It("should read and write files", func() {
		buf := &bytes.Buffer{}
		fs := tar.NewFs(tar.NewWriter(buf))

		file1, err := fs.Create("file1.txt")
		Expect(err).ToNot(HaveOccurred())
		Expect(file1.Close()).To(Succeed())

		r := tar.NewReader(buf)
		header, err := r.Next()
		Expect(err).ToNot(HaveOccurred())
		Expect(header.Name).To(Equal("file1.txt"))
		Expect(header.Size).To(Equal(int64(0)))

		file2, err := fs.Create("file2.txt")
		Expect(err).ToNot(HaveOccurred())
		Expect(file2.Close()).To(Succeed())

		header, err = r.Next()
		Expect(err).ToNot(HaveOccurred())
		Expect(header.Name).To(Equal("file2.txt"))
		Expect(header.Size).To(Equal(int64(0)))
	})
})
