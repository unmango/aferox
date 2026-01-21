package filter_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/spf13/afero"
	"github.com/unmango/aferox/filter"
	"github.com/unmango/aferox/testing"
)

var _ = Describe("Fs", func() {
	When("a single file exists", func() {
		var baseFs afero.Fs

		BeforeEach(func() {
			baseFs = afero.NewMemMapFs()
			err := afero.WriteFile(baseFs, "test.txt", []byte("testing"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
		})

		When("the file is filtered", func() {
			var filterFs afero.Fs

			BeforeEach(func() {
				filterFs = filter.NewFs(baseFs, func(name string) bool {
					return name == "not-test.txt"
				})
			})

			It("should not allow chmod-ing a filtered file", func() {
				err := filterFs.Chmod("test.txt", os.ModeAppend)

				Expect(err).To(MatchError(syscall.ENOENT))
			})

			It("should not allow chown-ing a filtered file", func() {
				err := filterFs.Chown("test.txt", 1001, 1001)

				Expect(err).To(MatchError(syscall.ENOENT))
			})

			It("should not allow chtimes-ing a filtered file", func() {
				err := filterFs.Chtimes("test.txt", time.Now(), time.Now())

				Expect(err).To(MatchError(syscall.ENOENT))
			})

			It("should not allow creating a filtered file", func() {
				_, err := filterFs.Create("test.txt")

				Expect(err).To(MatchError(syscall.ENOENT))
			})

			It("should not allow opening a filtered file", func() {
				_, err := filterFs.Open("test.txt")

				Expect(err).To(MatchError(syscall.ENOENT))
			})

			It("should not allow open-file-ing a filtered file", func() {
				_, err := filterFs.OpenFile("test.txt", 69, os.ModeAppend)

				Expect(err).To(MatchError(syscall.ENOENT))
			})

			It("should not allow removing a filtered file", func() {
				err := filterFs.Remove("test.txt")

				Expect(err).To(MatchError(syscall.ENOENT))
			})

			It("should not allow stat-ing a filtered file", func() {
				_, err := filterFs.Stat("test.txt")

				Expect(err).To(MatchError(syscall.ENOENT))
			})
		})

		When("the file is not filtered", func() {
			var filtered afero.Fs

			BeforeEach(func() {
				filtered = filter.NewFs(baseFs, func(name string) bool {
					return name == "test.txt"
				})
			})

			It("should allow chmod-ing a non-filtered file", func() {
				err := filtered.Chmod("test.txt", os.ModeAppend)

				Expect(err).NotTo(HaveOccurred())
			})

			It("should allow chown-ing a non-filtered file", func() {
				err := filtered.Chown("test.txt", 1001, 1001)

				Expect(err).NotTo(HaveOccurred())
			})

			It("should allow chtimes-ing a non-filtered file", func() {
				err := filtered.Chtimes("test.txt", time.Now(), time.Now())

				Expect(err).NotTo(HaveOccurred())
			})

			It("should allow creating a non-filtered file", func() {
				_, err := filtered.Create("test.txt")

				Expect(err).NotTo(HaveOccurred())
			})

			It("should include the source filesystem name", func() {
				Expect(filtered.Name()).To(Equal("Filter: MemMapFS"))
			})

			It("should allow opening a non-filtered file", func() {
				_, err := filtered.Open("test.txt")

				Expect(err).NotTo(HaveOccurred())
			})

			It("should allow open-file-ing a non-filtered file", func() {
				_, err := filtered.OpenFile("test.txt", 69, os.ModeAppend)

				Expect(err).NotTo(HaveOccurred())
			})

			It("should allow removing a non-filtered file", func() {
				err := filtered.Remove("test.txt")

				Expect(err).NotTo(HaveOccurred())
			})

			It("should allow removing a directory", func() {
				err := baseFs.Mkdir("test", os.ModeDir)
				Expect(err).NotTo(HaveOccurred())

				err = filtered.Remove("test")

				Expect(err).NotTo(HaveOccurred())
			})

			It("should allow stat-ing a non-filtered file", func() {
				_, err := filtered.Stat("test.txt")

				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	It("should walk properly", func() {
		base := afero.NewMemMapFs()
		Expect(base.Mkdir("test", os.ModePerm)).To(Succeed())
		err := afero.WriteFile(base, "test/file.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		filtered := filter.NewFs(base, func(s string) bool {
			return filepath.Ext(s) != ".txt"
		})
		paths := []string{}

		err = afero.Walk(filtered, "", func(path string, info fs.FileInfo, err error) error {
			paths = append(paths, path)
			return err
		})

		Expect(err).NotTo(HaveOccurred())
		Expect(paths).To(ConsistOf("", "test"))
	})

	It("should pass directory entries to the predicate", func() {
		base := afero.NewMemMapFs()
		Expect(base.Mkdir("test", os.ModePerm)).To(Succeed())
		err := afero.WriteFile(base, "test/file.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		filtered := filter.NewFs(base, func(s string) bool {
			Expect(s).To(Equal("test/file.txt"))
			return true
		})

		_, err = afero.ReadDir(filtered, "test")
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Mkdir", func() {
		It("should create a directory", func() {
			base := afero.NewMemMapFs()
			filtered := filter.NewFs(base, func(s string) bool { return true })

			err := filtered.Mkdir("testdir", os.ModePerm)

			Expect(err).NotTo(HaveOccurred())
			exists, err := afero.Exists(base, "testdir")
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())
		})
	})

	Describe("MkdirAll", func() {
		It("should create nested directories", func() {
			base := afero.NewMemMapFs()
			filtered := filter.NewFs(base, func(s string) bool { return true })

			err := filtered.MkdirAll("path/to/dir", os.ModePerm)

			Expect(err).NotTo(HaveOccurred())
			exists, err := afero.Exists(base, "path/to/dir")
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())
		})
	})

	Describe("RemoveAll", func() {
		It("should remove a directory", func() {
			base := afero.NewMemMapFs()
			Expect(base.Mkdir("testdir", os.ModePerm)).To(Succeed())
			filtered := filter.NewFs(base, func(s string) bool { return true })

			err := filtered.RemoveAll("testdir")

			Expect(err).NotTo(HaveOccurred())
			exists, err := afero.Exists(base, "testdir")
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())
		})

		It("should remove a non-filtered file", func() {
			base := afero.NewMemMapFs()
			Expect(afero.WriteFile(base, "test.txt", []byte("test"), os.ModePerm)).To(Succeed())
			filtered := filter.NewFs(base, func(s string) bool { return s == "test.txt" })

			err := filtered.RemoveAll("test.txt")

			Expect(err).NotTo(HaveOccurred())
			exists, err := afero.Exists(base, "test.txt")
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())
		})

		It("should not allow removing a filtered file", func() {
			base := afero.NewMemMapFs()
			Expect(afero.WriteFile(base, "test.txt", []byte("test"), os.ModePerm)).To(Succeed())
			filtered := filter.NewFs(base, func(s string) bool { return s != "test.txt" })

			err := filtered.RemoveAll("test.txt")

			Expect(err).To(MatchError(syscall.ENOENT))
		})
	})

	Describe("Rename", func() {
		It("should rename a non-filtered file", func() {
			base := afero.NewMemMapFs()
			Expect(afero.WriteFile(base, "old.txt", []byte("test"), os.ModePerm)).To(Succeed())
			filtered := filter.NewFs(base, func(s string) bool { return s == "old.txt" || s == "new.txt" })

			err := filtered.Rename("old.txt", "new.txt")

			Expect(err).NotTo(HaveOccurred())
			exists, err := afero.Exists(base, "new.txt")
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())
		})

		It("should not rename a filtered file", func() {
			base := afero.NewMemMapFs()
			Expect(afero.WriteFile(base, "old.txt", []byte("test"), os.ModePerm)).To(Succeed())
			filtered := filter.NewFs(base, func(s string) bool { return false })

			err := filtered.Rename("old.txt", "new.txt")

			Expect(err).To(MatchError(syscall.ENOENT))
		})

		It("should not rename to a filtered name", func() {
			base := afero.NewMemMapFs()
			Expect(afero.WriteFile(base, "old.txt", []byte("test"), os.ModePerm)).To(Succeed())
			filtered := filter.NewFs(base, func(s string) bool { return s == "old.txt" })

			err := filtered.Rename("old.txt", "new.txt")

			Expect(err).To(MatchError(syscall.ENOENT))
		})

		It("should allow renaming directories", func() {
			base := afero.NewMemMapFs()
			Expect(base.Mkdir("olddir", os.ModePerm)).To(Succeed())
			filtered := filter.NewFs(base, func(s string) bool { return false })

			err := filtered.Rename("olddir", "newdir")

			Expect(err).NotTo(HaveOccurred())
		})
	})

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

	Describe("Filter variations", func() {
		It("should work with nil filter", func() {
			base := afero.NewMemMapFs()
			Expect(afero.WriteFile(base, "test.txt", []byte("test"), os.ModePerm)).To(Succeed())
			filtered := filter.FromFilter(base, nil)

			_, err := filtered.Open("test.txt")
			Expect(err).NotTo(HaveOccurred())

			_, err = filtered.Create("new.txt")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should work with FromPredicateWithError", func() {
			base := afero.NewMemMapFs()
			Expect(afero.WriteFile(base, "test.txt", []byte("test"), os.ModePerm)).To(Succeed())
			customErr := fs.ErrPermission
			filtered := filter.FromPredicateWithError(base, func(s string) bool {
				return s != "test.txt"
			}, customErr)

			_, err := filtered.Open("test.txt")
			Expect(err).To(MatchError(customErr))
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

	Describe("Error handling", func() {
		It("should handle IsDir error in Open", func() {
			base := &testing.Fs{Fs: afero.NewMemMapFs()}
			expectedErr := fs.ErrPermission
			base.StatFunc = func(name string) (fs.FileInfo, error) {
				return nil, expectedErr
			}
			
			filtered := filter.NewFs(base, func(s string) bool { return true })
			_, err := filtered.Open("test.txt")
			
			Expect(err).To(MatchError(expectedErr))
		})

		It("should handle Open error", func() {
			base := &testing.Fs{Fs: afero.NewMemMapFs()}
			Expect(afero.WriteFile(base, "test.txt", []byte("test"), os.ModePerm)).To(Succeed())
			
			expectedErr := fs.ErrNotExist
			base.OpenFunc = func(name string) (afero.File, error) {
				return nil, expectedErr
			}
			
			filtered := filter.NewFs(base, func(s string) bool { return true })
			_, err := filtered.Open("test.txt")
			
			Expect(err).To(MatchError(expectedErr))
		})

		It("should handle IsDir error in RemoveAll", func() {
			base := &testing.Fs{Fs: afero.NewMemMapFs()}
			expectedErr := fs.ErrPermission
			base.StatFunc = func(name string) (fs.FileInfo, error) {
				return nil, expectedErr
			}
			
			filtered := filter.NewFs(base, func(s string) bool { return true })
			err := filtered.RemoveAll("test.txt")
			
			Expect(err).To(MatchError(expectedErr))
		})

		It("should handle IsDir error in Rename", func() {
			base := &testing.Fs{Fs: afero.NewMemMapFs()}
			expectedErr := fs.ErrPermission
			base.StatFunc = func(name string) (fs.FileInfo, error) {
				return nil, expectedErr
			}
			
			filtered := filter.NewFs(base, func(s string) bool { return true })
			err := filtered.Rename("old.txt", "new.txt")
			
			Expect(err).To(MatchError(expectedErr))
		})

		It("should handle IsDir error in dirOrMatches", func() {
			base := &testing.Fs{Fs: afero.NewMemMapFs()}
			expectedErr := fs.ErrInvalid
			base.StatFunc = func(name string) (fs.FileInfo, error) {
				return nil, expectedErr
			}
			
			filtered := filter.NewFs(base, func(s string) bool { return true })
			
			// Test dirOrMatches through Stat
			_, err := filtered.Stat("test.txt")
			Expect(err).To(MatchError(expectedErr))
		})
	})
})
