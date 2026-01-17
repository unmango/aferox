package tar_test

import (
	"bytes"
	"syscall"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	aferoxtar "github.com/unmango/aferox/writer/tar"
)

var _ = Describe("File", func() {
	It("should support Write method", func() {
		buf := &bytes.Buffer{}
		fs := aferoxtar.NewFs(aferoxtar.NewWriter(buf))

		file, err := fs.Create("test.txt")
		Expect(err).ToNot(HaveOccurred())

		n, err := file.Write([]byte("Hello"))
		Expect(err).ToNot(HaveOccurred())
		Expect(n).To(Equal(5))

		Expect(file.Close()).To(Succeed())

		r := aferoxtar.NewReader(buf)
		header, err := r.Next()
		Expect(err).ToNot(HaveOccurred())
		Expect(header.Name).To(Equal("test.txt"))
		Expect(header.Size).To(Equal(int64(5)))
	})

	It("should return file name", func() {
		buf := &bytes.Buffer{}
		fs := aferoxtar.NewFs(aferoxtar.NewWriter(buf))

		file, err := fs.Create("myfile.txt")
		Expect(err).ToNot(HaveOccurred())
		Expect(file.Name()).To(Equal("myfile.txt"))
		Expect(file.Close()).To(Succeed())
	})

	It("should return EROFS for Read operations", func() {
		buf := &bytes.Buffer{}
		fs := aferoxtar.NewFs(aferoxtar.NewWriter(buf))

		file, err := fs.Create("test.txt")
		Expect(err).ToNot(HaveOccurred())

		data := make([]byte, 10)
		_, err = file.Read(data)
		Expect(err).To(Equal(syscall.EROFS))

		Expect(file.Close()).To(Succeed())
	})

	It("should return EROFS for ReadAt operations", func() {
		buf := &bytes.Buffer{}
		fs := aferoxtar.NewFs(aferoxtar.NewWriter(buf))

		file, err := fs.Create("test.txt")
		Expect(err).ToNot(HaveOccurred())

		data := make([]byte, 10)
		_, err = file.ReadAt(data, 0)
		Expect(err).To(Equal(syscall.EROFS))

		Expect(file.Close()).To(Succeed())
	})

	It("should return EROFS for Readdir operations", func() {
		buf := &bytes.Buffer{}
		fs := aferoxtar.NewFs(aferoxtar.NewWriter(buf))

		file, err := fs.Create("test.txt")
		Expect(err).ToNot(HaveOccurred())

		_, err = file.Readdir(0)
		Expect(err).To(Equal(syscall.EROFS))

		Expect(file.Close()).To(Succeed())
	})

	It("should return EROFS for Readdirnames operations", func() {
		buf := &bytes.Buffer{}
		fs := aferoxtar.NewFs(aferoxtar.NewWriter(buf))

		file, err := fs.Create("test.txt")
		Expect(err).ToNot(HaveOccurred())

		_, err = file.Readdirnames(0)
		Expect(err).To(Equal(syscall.EROFS))

		Expect(file.Close()).To(Succeed())
	})

	It("should return EROFS for Seek operations", func() {
		buf := &bytes.Buffer{}
		fs := aferoxtar.NewFs(aferoxtar.NewWriter(buf))

		file, err := fs.Create("test.txt")
		Expect(err).ToNot(HaveOccurred())

		_, err = file.Seek(0, 0)
		Expect(err).To(Equal(syscall.EROFS))

		Expect(file.Close()).To(Succeed())
	})

	It("should return FileInfo from Stat", func() {
		buf := &bytes.Buffer{}
		fs := aferoxtar.NewFs(aferoxtar.NewWriter(buf))

		file, err := fs.Create("test.txt")
		Expect(err).ToNot(HaveOccurred())

		info, err := file.Stat()
		Expect(err).ToNot(HaveOccurred())
		Expect(info).ToNot(BeNil())
		Expect(info.Name()).To(Equal("test.txt"))

		Expect(file.Close()).To(Succeed())
	})

	It("should return EROFS for Sync operations", func() {
		buf := &bytes.Buffer{}
		fs := aferoxtar.NewFs(aferoxtar.NewWriter(buf))

		file, err := fs.Create("test.txt")
		Expect(err).ToNot(HaveOccurred())

		err = file.Sync()
		Expect(err).To(Equal(syscall.EROFS))

		Expect(file.Close()).To(Succeed())
	})

	It("should return EROFS for Truncate operations", func() {
		buf := &bytes.Buffer{}
		fs := aferoxtar.NewFs(aferoxtar.NewWriter(buf))

		file, err := fs.Create("test.txt")
		Expect(err).ToNot(HaveOccurred())

		err = file.Truncate(0)
		Expect(err).To(Equal(syscall.EROFS))

		Expect(file.Close()).To(Succeed())
	})

	It("should return EROFS for WriteAt operations", func() {
		buf := &bytes.Buffer{}
		fs := aferoxtar.NewFs(aferoxtar.NewWriter(buf))

		file, err := fs.Create("test.txt")
		Expect(err).ToNot(HaveOccurred())

		_, err = file.WriteAt([]byte("test"), 0)
		Expect(err).To(Equal(syscall.EROFS))

		Expect(file.Close()).To(Succeed())
	})
})
