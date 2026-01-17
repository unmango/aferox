package tar_test

import (
	"archive/tar"
	"bytes"
	"errors"
	"os"
	"syscall"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	aferoxtar "github.com/unmango/aferox/writer/tar"
)

// failingWriter is a writer that always fails
type failingWriter struct{}

func (w *failingWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("write failed")
}

// conditionalWriter succeeds for small writes (header) but fails for larger writes (content)
type conditionalWriter struct {
	maxSize int
}

func (w *conditionalWriter) Write(p []byte) (n int, err error) {
	// tar header is 512 bytes, so accept anything <= 512 bytes
	// and fail on larger writes (which would be content)
	if len(p) > w.maxSize {
		return 0, errors.New("write failed during content copy")
	}
	return len(p), nil
}

var _ = Describe("Fs", func() {
	It("should create a file", func() {
		buf := &bytes.Buffer{}
		fs := aferoxtar.NewFs(aferoxtar.NewWriter(buf))

		file, err := fs.Create("test.txt")
		Expect(err).ToNot(HaveOccurred())
		defer file.Close()

		_, err = file.WriteString("Hello, World!")
		Expect(err).ToNot(HaveOccurred())
		Expect(file.Close()).To(Succeed())

		r := aferoxtar.NewReader(buf)
		header, err := r.Next()
		Expect(err).ToNot(HaveOccurred())
		Expect(header.Name).To(Equal("test.txt"))
		Expect(header.Size).To(Equal(int64(13)))
	})

	It("should write multiple files", func() {
		buf := &bytes.Buffer{}
		fs := aferoxtar.NewFs(aferoxtar.NewWriter(buf))

		file1, err := fs.Create("file1.txt")
		Expect(err).ToNot(HaveOccurred())
		Expect(file1.Close()).To(Succeed())

		file2, err := fs.Create("file2.txt")
		Expect(err).ToNot(HaveOccurred())
		Expect(file2.Close()).To(Succeed())

		file3, err := fs.Create("file3.txt")
		Expect(err).ToNot(HaveOccurred())
		Expect(file3.Close()).To(Succeed())

		r := aferoxtar.NewReader(buf)
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
		fs := aferoxtar.NewFs(aferoxtar.NewWriter(buf))

		file1, err := fs.Create("file1.txt")
		Expect(err).ToNot(HaveOccurred())
		Expect(file1.Close()).To(Succeed())

		r := aferoxtar.NewReader(buf)
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

	Describe("File methods", func() {
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

	Describe("FileInfo methods", func() {
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

	Describe("Fs methods", func() {
		It("should return EPERM for Chmod", func() {
			buf := &bytes.Buffer{}
			fs := aferoxtar.NewFs(aferoxtar.NewWriter(buf))

			err := fs.Chmod("test.txt", 0644)
			Expect(err).To(Equal(syscall.EPERM))
		})

		It("should return EPERM for Chown", func() {
			buf := &bytes.Buffer{}
			fs := aferoxtar.NewFs(aferoxtar.NewWriter(buf))

			err := fs.Chown("test.txt", 0, 0)
			Expect(err).To(Equal(syscall.EPERM))
		})

		It("should return EPERM for Chtimes", func() {
			buf := &bytes.Buffer{}
			fs := aferoxtar.NewFs(aferoxtar.NewWriter(buf))

			err := fs.Chtimes("test.txt", time.Now(), time.Now())
			Expect(err).To(Equal(syscall.EPERM))
		})

		It("should create directories", func() {
			buf := &bytes.Buffer{}
			tw := aferoxtar.NewWriter(buf)
			fs := aferoxtar.NewFs(tw)

			err := fs.Mkdir("testdir", 0755)
			Expect(err).ToNot(HaveOccurred())

			tw.Close()

			r := aferoxtar.NewReader(buf)
			header, err := r.Next()
			Expect(err).ToNot(HaveOccurred())
			Expect(header.Name).To(Equal("testdir"))
			Expect(header.Typeflag).To(Equal(byte(tar.TypeDir)))
		})

		It("should create directories with MkdirAll", func() {
			buf := &bytes.Buffer{}
			tw := aferoxtar.NewWriter(buf)
			fs := aferoxtar.NewFs(tw)

			err := fs.MkdirAll("testdir/subdir", 0755)
			Expect(err).ToNot(HaveOccurred())

			tw.Close()

			r := aferoxtar.NewReader(buf)
			header, err := r.Next()
			Expect(err).ToNot(HaveOccurred())
			Expect(header.Name).To(Equal("testdir/subdir"))
			Expect(header.Typeflag).To(Equal(byte(tar.TypeDir)))
		})

		It("should return filesystem name", func() {
			buf := &bytes.Buffer{}
			fs := aferoxtar.NewFs(aferoxtar.NewWriter(buf))

			Expect(fs.Name()).To(Equal("tar.Writer"))
		})

		It("should return EROFS for Open", func() {
			buf := &bytes.Buffer{}
			fs := aferoxtar.NewFs(aferoxtar.NewWriter(buf))

			_, err := fs.Open("test.txt")
			Expect(err).To(Equal(syscall.EROFS))
		})

		It("should support OpenFile", func() {
			buf := &bytes.Buffer{}
			tw := aferoxtar.NewWriter(buf)
			fs := aferoxtar.NewFs(tw)

			file, err := fs.OpenFile("test.txt", 0, 0644)
			Expect(err).ToNot(HaveOccurred())
			Expect(file).ToNot(BeNil())

			_, err = file.Write([]byte("content"))
			Expect(err).ToNot(HaveOccurred())
			Expect(file.Close()).To(Succeed())

			tw.Close()

			r := aferoxtar.NewReader(buf)
			header, err := r.Next()
			Expect(err).ToNot(HaveOccurred())
			Expect(header.Name).To(Equal("test.txt"))
			Expect(header.Size).To(Equal(int64(7)))
		})

		It("should return EROFS for Remove", func() {
			buf := &bytes.Buffer{}
			fs := aferoxtar.NewFs(aferoxtar.NewWriter(buf))

			err := fs.Remove("test.txt")
			Expect(err).To(Equal(syscall.EROFS))
		})

		It("should return EROFS for RemoveAll", func() {
			buf := &bytes.Buffer{}
			fs := aferoxtar.NewFs(aferoxtar.NewWriter(buf))

			err := fs.RemoveAll("test.txt")
			Expect(err).To(Equal(syscall.EROFS))
		})

		It("should return EROFS for Rename", func() {
			buf := &bytes.Buffer{}
			fs := aferoxtar.NewFs(aferoxtar.NewWriter(buf))

			err := fs.Rename("old.txt", "new.txt")
			Expect(err).To(Equal(syscall.EROFS))
		})

		It("should return FileInfo for Stat", func() {
			buf := &bytes.Buffer{}
			fs := aferoxtar.NewFs(aferoxtar.NewWriter(buf))

			info, err := fs.Stat("test.txt")
			Expect(err).ToNot(HaveOccurred())
			Expect(info).ToNot(BeNil())
			Expect(info.Name()).To(Equal("test.txt"))
		})
	})

	Describe("Error handling", func() {
		It("should handle flush errors when writing header fails", func() {
			fw := &failingWriter{}
			tw := tar.NewWriter(fw)
			fs := aferoxtar.NewFs(tw)

			file, err := fs.Create("test.txt")
			Expect(err).ToNot(HaveOccurred())

			_, err = file.Write([]byte("content"))
			Expect(err).ToNot(HaveOccurred())

			err = file.Close()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("writing header"))
		})

		It("should handle flush errors when copying buffer fails", func() {
			// Create a tar writer that accepts the header but fails on copy
			cw := &conditionalWriter{maxSize: 512}
			tw := tar.NewWriter(cw)
			fs := aferoxtar.NewFs(tw)

			file, err := fs.Create("test.txt")
			Expect(err).ToNot(HaveOccurred())

			// Write content larger than 512 bytes to trigger failure during copy
			_, err = file.Write(bytes.Repeat([]byte("x"), 1024))
			Expect(err).ToNot(HaveOccurred())

			err = file.Close()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("copying file buffer"))
		})

		It("should handle Mkdir errors when writing header fails", func() {
			fw := &failingWriter{}
			tw := tar.NewWriter(fw)
			fs := aferoxtar.NewFs(tw)

			err := fs.Mkdir("testdir", 0755)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("writing header"))
		})
	})
})
