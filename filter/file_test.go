package filter_test

import (
	"io/fs"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/spf13/afero"
	"github.com/unmango/aferox/filter"
	"github.com/unmango/aferox/testing"
)

var _ = Describe("File", func() {
	Describe("File operations", func() {
		var base afero.Fs
		var filtered afero.Fs
		var file afero.File

		BeforeEach(func() {
			base = afero.NewMemMapFs()
			Expect(afero.WriteFile(base, "test.txt", []byte("hello world"), os.ModePerm)).To(Succeed())
			filtered = filter.NewFs(base, func(s string) bool { return s == "test.txt" })
			var err error
			file, err = filtered.Open("test.txt")
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			if file != nil {
				file.Close()
			}
		})

		It("should return file name", func() {
			Expect(file.Name()).To(Equal("test.txt"))
		})

		It("should read from file", func() {
			buf := make([]byte, 5)
			n, err := file.Read(buf)
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(5))
			Expect(string(buf)).To(Equal("hello"))
		})

		It("should read at offset", func() {
			buf := make([]byte, 5)
			n, err := file.ReadAt(buf, 6)
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(5))
			Expect(string(buf)).To(Equal("world"))
		})

		It("should seek in file", func() {
			pos, err := file.Seek(6, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(pos).To(Equal(int64(6)))
		})

		It("should stat file", func() {
			info, err := file.Stat()
			Expect(err).NotTo(HaveOccurred())
			Expect(info.Name()).To(Equal("test.txt"))
			Expect(info.Size()).To(Equal(int64(11)))
		})

		It("should sync file", func() {
			err := file.Sync()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should truncate writable file", func() {
			file.Close()
			wf, err := filtered.OpenFile("test.txt", os.O_RDWR, os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
			defer wf.Close()

			err = wf.Truncate(5)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should write to writable file", func() {
			file.Close()
			wf, err := filtered.OpenFile("test.txt", os.O_RDWR, os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
			defer wf.Close()

			n, err := wf.Write([]byte(" test"))
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(5))
		})

		It("should write at offset in writable file", func() {
			file.Close()
			wf, err := filtered.OpenFile("test.txt", os.O_RDWR, os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
			defer wf.Close()

			n, err := wf.WriteAt([]byte("TEST"), 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(4))
		})

		It("should write string to writable file", func() {
			file.Close()
			wf, err := filtered.OpenFile("test.txt", os.O_RDWR, os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
			defer wf.Close()

			n, err := wf.WriteString(" testing")
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(8))
		})

		Describe("Readdir", func() {
			BeforeEach(func() {
				if file != nil {
					file.Close()
				}
				Expect(base.Mkdir("dir", os.ModePerm)).To(Succeed())
				Expect(afero.WriteFile(base, "dir/file1.txt", []byte("1"), os.ModePerm)).To(Succeed())
				Expect(afero.WriteFile(base, "dir/file2.go", []byte("2"), os.ModePerm)).To(Succeed())
				Expect(base.Mkdir("dir/subdir", os.ModePerm)).To(Succeed())

				filtered = filter.NewFs(base, func(s string) bool {
					return filepath.Ext(s) == ".txt"
				})

				var err error
				file, err = filtered.Open("dir")
				Expect(err).NotTo(HaveOccurred())
			})

			It("should filter files in Readdir", func() {
				infos, err := file.Readdir(-1)
				Expect(err).NotTo(HaveOccurred())
				Expect(infos).To(HaveLen(2))
				names := []string{infos[0].Name(), infos[1].Name()}
				Expect(names).To(ConsistOf("file1.txt", "subdir"))
			})
		})

		Describe("Readdirnames", func() {
			BeforeEach(func() {
				if file != nil {
					file.Close()
				}
				Expect(base.Mkdir("dir", os.ModePerm)).To(Succeed())
				Expect(afero.WriteFile(base, "dir/file1.txt", []byte("1"), os.ModePerm)).To(Succeed())
				Expect(afero.WriteFile(base, "dir/file2.go", []byte("2"), os.ModePerm)).To(Succeed())

				filtered = filter.NewFs(base, func(s string) bool {
					return filepath.Ext(s) == ".txt"
				})

				var err error
				file, err = filtered.Open("dir")
				Expect(err).NotTo(HaveOccurred())
			})

			It("should filter files in Readdirnames", func() {
				names, err := file.Readdirnames(-1)
				Expect(err).NotTo(HaveOccurred())
				Expect(names).To(ConsistOf("file1.txt"))
			})
		})
	})

	Describe("Direct File wrapper operations", func() {
		It("should handle write operations on wrapped file", func() {
			base := &testing.Fs{Fs: afero.NewMemMapFs()}
			Expect(afero.WriteFile(base, "test.txt", []byte("hello"), os.ModePerm)).To(Succeed())

			// Mock Open to return a writable file
			base.OpenFunc = func(name string) (afero.File, error) {
				return base.Fs.OpenFile(name, os.O_RDWR, os.ModePerm)
			}

			filtered := filter.NewFs(base, func(s string) bool { return s == "test.txt" })

			// Now Open will return a File wrapper around a writable file
			file, err := filtered.Open("test.txt")
			Expect(err).NotTo(HaveOccurred())
			defer file.Close()

			// Test write operations
			n, err := file.Write([]byte(" world"))
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(6))

			n, err = file.WriteAt([]byte("HELLO"), 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(5))

			n, err = file.WriteString("!")
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(1))

			err = file.Truncate(3)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should handle Readdir and Readdirnames errors", func() {
			base := &testing.Fs{Fs: afero.NewMemMapFs()}
			Expect(base.Mkdir("dir", os.ModePerm)).To(Succeed())
			Expect(afero.WriteFile(base, "dir/file.txt", []byte("test"), os.ModePerm)).To(Succeed())

			errToReturn := fs.ErrInvalid
			base.OpenFunc = func(name string) (afero.File, error) {
				return &testing.File{
					CloseFunc: func() error { return nil },
					ReaddirFunc: func(count int) ([]fs.FileInfo, error) {
						return nil, errToReturn
					},
					ReaddirnamesFunc: func(n int) ([]string, error) {
						return nil, errToReturn
					},
				}, nil
			}

			filtered := filter.NewFs(base, func(s string) bool { return true })
			f, err := filtered.Open("dir")
			Expect(err).NotTo(HaveOccurred())
			defer f.Close()

			// These should propagate the error
			_, err = f.Readdir(-1)
			Expect(err).To(MatchError(errToReturn))

			_, err = f.Readdirnames(-1)
			Expect(err).To(MatchError(errToReturn))
		})
	})
})
