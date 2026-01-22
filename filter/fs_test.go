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
	"github.com/unmango/aferox/op"
	"github.com/unmango/aferox/testing"
)

var _ = Describe("Fs", func() {
	When("a single file exists", func() {
		var baseFs afero.Fs

		BeforeEach(func() {
			baseFs = afero.NewMemMapFs()
			err := afero.WriteFile(baseFs,
				"test.txt",
				[]byte("testing"),
				os.ModePerm,
			)
			Expect(err).NotTo(HaveOccurred())
		})

		When("the file is filtered", func() {
			var filterFs afero.Fs

			BeforeEach(func() {
				filterFs = filter.FromPredicate(baseFs, func(o op.Operation) bool {
					return o.Path() == "not-test.txt"
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
				filtered = filter.FromPredicate(baseFs, func(o op.Operation) bool {
					return o.Path() == "test.txt"
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
		filtered := filter.FromPredicate(base, func(o op.Operation) bool {
			return filepath.Ext(o.Path()) != ".txt"
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

		filtered := filter.FromPredicate(base, func(o op.Operation) bool {
			Expect(o.Path()).To(Equal("test/file.txt"))
			return true
		})

		_, err = afero.ReadDir(filtered, "test")
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Mkdir", func() {
		It("should create a directory", func() {
			base := afero.NewMemMapFs()
			filtered := filter.FromPredicate(base, func(o op.Operation) bool {
				return true
			})

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
			filtered := filter.FromPredicate(base, func(o op.Operation) bool {
				return true
			})

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
			filtered := filter.FromPredicate(base, func(o op.Operation) bool {
				return true
			})

			err := filtered.RemoveAll("testdir")

			Expect(err).NotTo(HaveOccurred())
			exists, err := afero.Exists(base, "testdir")
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())
		})

		It("should remove a non-filtered file", func() {
			base := afero.NewMemMapFs()
			Expect(afero.WriteFile(base, "test.txt", []byte("test"), os.ModePerm)).To(Succeed())
			filtered := filter.FromPredicate(base, func(o op.Operation) bool {
				return o.Path() == "test.txt"
			})

			err := filtered.RemoveAll("test.txt")

			Expect(err).NotTo(HaveOccurred())
			exists, err := afero.Exists(base, "test.txt")
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())
		})

		It("should not allow removing a filtered file", func() {
			base := afero.NewMemMapFs()
			Expect(afero.WriteFile(base, "test.txt", []byte("test"), os.ModePerm)).To(Succeed())
			filtered := filter.FromPredicate(base, func(o op.Operation) bool {
				return o.Path() != "test.txt"
			})

			err := filtered.RemoveAll("test.txt")

			Expect(err).To(MatchError(syscall.ENOENT))
		})
	})

	Describe("Rename", func() {
		It("should rename a non-filtered file", func() {
			base := afero.NewMemMapFs()
			Expect(afero.WriteFile(base, "old.txt", []byte("test"), os.ModePerm)).To(Succeed())
			filtered := filter.FromPredicate(base, func(o op.Operation) bool {
				return o.Path() == "old.txt" || o.Path() == "new.txt"
			})

			err := filtered.Rename("old.txt", "new.txt")

			Expect(err).NotTo(HaveOccurred())
			exists, err := afero.Exists(base, "new.txt")
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())
		})

		It("should not rename a filtered file", func() {
			base := afero.NewMemMapFs()
			Expect(afero.WriteFile(base, "old.txt", []byte("test"), os.ModePerm)).To(Succeed())
			filtered := filter.FromPredicate(base, func(o op.Operation) bool {
				return false
			})

			err := filtered.Rename("old.txt", "new.txt")

			Expect(err).To(MatchError(syscall.ENOENT))
		})

		It("should not rename to a filtered name", func() {
			base := afero.NewMemMapFs()
			Expect(afero.WriteFile(base, "old.txt", []byte("test"), os.ModePerm)).To(Succeed())
			filtered := filter.FromPredicate(base, func(o op.Operation) bool {
				return o.Path() == "old.txt"
			})

			err := filtered.Rename("old.txt", "new.txt")

			Expect(err).To(MatchError(syscall.ENOENT))
		})

		It("should allow renaming directories", func() {
			base := afero.NewMemMapFs()
			Expect(base.Mkdir("olddir", os.ModePerm)).To(Succeed())
			filtered := filter.FromPredicate(base, func(o op.Operation) bool {
				return false
			})

			err := filtered.Rename("olddir", "newdir")

			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Error handling", func() {
		It("should handle IsDir error in Open", func() {
			base := &testing.Fs{Fs: afero.NewMemMapFs()}
			expectedErr := fs.ErrPermission
			base.StatFunc = func(name string) (fs.FileInfo, error) {
				return nil, expectedErr
			}

			filtered := filter.FromPredicate(base, func(o op.Operation) bool {
				return true
			})
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

			filtered := filter.FromPredicate(base, func(o op.Operation) bool {
				return true
			})
			_, err := filtered.Open("test.txt")

			Expect(err).To(MatchError(expectedErr))
		})

		It("should handle IsDir error in RemoveAll", func() {
			base := &testing.Fs{Fs: afero.NewMemMapFs()}
			expectedErr := fs.ErrPermission
			base.StatFunc = func(name string) (fs.FileInfo, error) {
				return nil, expectedErr
			}

			filtered := filter.FromPredicate(base, func(o op.Operation) bool {
				return true
			})
			err := filtered.RemoveAll("test.txt")

			Expect(err).To(MatchError(expectedErr))
		})

		It("should handle IsDir error in Rename", func() {
			base := &testing.Fs{Fs: afero.NewMemMapFs()}
			expectedErr := fs.ErrPermission
			base.StatFunc = func(name string) (fs.FileInfo, error) {
				return nil, expectedErr
			}

			filtered := filter.FromPredicate(base, func(o op.Operation) bool {
				return true
			})
			err := filtered.Rename("old.txt", "new.txt")

			Expect(err).To(MatchError(expectedErr))
		})

		It("should handle IsDir error in dirOrMatches", func() {
			base := &testing.Fs{Fs: afero.NewMemMapFs()}
			expectedErr := fs.ErrInvalid
			base.StatFunc = func(name string) (fs.FileInfo, error) {
				return nil, expectedErr
			}

			filtered := filter.FromPredicate(base, func(o op.Operation) bool {
				return true
			})

			// Test dirOrMatches through Stat
			_, err := filtered.Stat("test.txt")
			Expect(err).To(MatchError(expectedErr))
		})
	})
})
