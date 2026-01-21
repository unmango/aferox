package filter_test

import (
	"io/fs"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/spf13/afero"
	"github.com/unmango/aferox/filter"
)

var _ = Describe("Operation", func() {
	Describe("Operation types", func() {
		It("should handle ChmodOp with type switch", func() {
			base := afero.NewMemMapFs()
			Expect(afero.WriteFile(base, "test.txt", []byte("test"), os.ModePerm)).To(Succeed())

			filtered := filter.NewFs(base, func(op filter.Operation) bool {
				switch v := op.(type) {
				case filter.ChmodOp:
					Expect(v.Name).To(Equal("test.txt"))
					Expect(v.Mode).To(Equal(fs.FileMode(0755)))
					Expect(v.Path()).To(Equal("test.txt"))
					return true
				default:
					return true
				}
			})

			err := filtered.Chmod("test.txt", 0755)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should handle ChownOp with type switch", func() {
			base := afero.NewMemMapFs()
			Expect(afero.WriteFile(base, "test.txt", []byte("test"), os.ModePerm)).To(Succeed())

			filtered := filter.NewFs(base, func(op filter.Operation) bool {
				switch v := op.(type) {
				case filter.ChownOp:
					Expect(v.Name).To(Equal("test.txt"))
					Expect(v.UID).To(Equal(1000))
					Expect(v.GID).To(Equal(1000))
					Expect(v.Path()).To(Equal("test.txt"))
					return true
				default:
					return true
				}
			})

			err := filtered.Chown("test.txt", 1000, 1000)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should handle ChtimesOp with type switch", func() {
			base := afero.NewMemMapFs()
			Expect(afero.WriteFile(base, "test.txt", []byte("test"), os.ModePerm)).To(Succeed())

			now := time.Now()
			filtered := filter.NewFs(base, func(op filter.Operation) bool {
				switch v := op.(type) {
				case filter.ChtimesOp:
					Expect(v.Name).To(Equal("test.txt"))
					Expect(v.Path()).To(Equal("test.txt"))
					return true
				default:
					return true
				}
			})

			err := filtered.Chtimes("test.txt", now, now)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should handle CreateOp with type switch", func() {
			base := afero.NewMemMapFs()

			filtered := filter.NewFs(base, func(op filter.Operation) bool {
				switch v := op.(type) {
				case filter.CreateOp:
					Expect(v.Name).To(Equal("test.txt"))
					Expect(v.Path()).To(Equal("test.txt"))
					return true
				default:
					return true
				}
			})

			_, err := filtered.Create("test.txt")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should handle MkdirOp with type switch", func() {
			base := afero.NewMemMapFs()

			filtered := filter.NewFs(base, func(op filter.Operation) bool {
				switch v := op.(type) {
				case filter.MkdirOp:
					Expect(v.Name).To(Equal("testdir"))
					Expect(v.Perm).To(Equal(os.ModePerm))
					Expect(v.Path()).To(Equal("testdir"))
					return true
				default:
					return true
				}
			})

			err := filtered.Mkdir("testdir", os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should handle MkdirAllOp with type switch", func() {
			base := afero.NewMemMapFs()

			filtered := filter.NewFs(base, func(op filter.Operation) bool {
				switch v := op.(type) {
				case filter.MkdirAllOp:
					Expect(v.PathName).To(Equal("path/to/dir"))
					Expect(v.Perm).To(Equal(os.ModePerm))
					Expect(v.Path()).To(Equal("path/to/dir"))
					return true
				default:
					return true
				}
			})

			err := filtered.MkdirAll("path/to/dir", os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should handle OpenOp with type switch", func() {
			base := afero.NewMemMapFs()
			Expect(afero.WriteFile(base, "test.txt", []byte("test"), os.ModePerm)).To(Succeed())

			filtered := filter.NewFs(base, func(op filter.Operation) bool {
				switch v := op.(type) {
				case filter.OpenOp:
					Expect(v.Name).To(Equal("test.txt"))
					Expect(v.Path()).To(Equal("test.txt"))
					return true
				default:
					return true
				}
			})

			_, err := filtered.Open("test.txt")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should handle OpenFileOp with type switch", func() {
			base := afero.NewMemMapFs()
			Expect(afero.WriteFile(base, "test.txt", []byte("test"), os.ModePerm)).To(Succeed())

			filtered := filter.NewFs(base, func(op filter.Operation) bool {
				switch v := op.(type) {
				case filter.OpenFileOp:
					Expect(v.Name).To(Equal("test.txt"))
					Expect(v.Flag).To(Equal(os.O_RDWR))
					Expect(v.Perm).To(Equal(os.ModePerm))
					Expect(v.Path()).To(Equal("test.txt"))
					return true
				default:
					return true
				}
			})

			_, err := filtered.OpenFile("test.txt", os.O_RDWR, os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should handle RemoveOp with type switch", func() {
			base := afero.NewMemMapFs()
			Expect(afero.WriteFile(base, "test.txt", []byte("test"), os.ModePerm)).To(Succeed())

			filtered := filter.NewFs(base, func(op filter.Operation) bool {
				switch v := op.(type) {
				case filter.RemoveOp:
					Expect(v.Name).To(Equal("test.txt"))
					Expect(v.Path()).To(Equal("test.txt"))
					return true
				default:
					return true
				}
			})

			err := filtered.Remove("test.txt")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should handle RemoveAllOp with type switch", func() {
			base := afero.NewMemMapFs()
			Expect(afero.WriteFile(base, "test.txt", []byte("test"), os.ModePerm)).To(Succeed())

			filtered := filter.NewFs(base, func(op filter.Operation) bool {
				switch v := op.(type) {
				case filter.RemoveAllOp:
					Expect(v.PathName).To(Equal("test.txt"))
					Expect(v.Path()).To(Equal("test.txt"))
					return true
				default:
					return true
				}
			})

			err := filtered.RemoveAll("test.txt")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should handle RenameOp with type switch", func() {
			base := afero.NewMemMapFs()
			Expect(afero.WriteFile(base, "old.txt", []byte("test"), os.ModePerm)).To(Succeed())

			renameSeen := false
			createSeen := false

			filtered := filter.NewFs(base, func(op filter.Operation) bool {
				switch v := op.(type) {
				case filter.RenameOp:
					Expect(v.Oldname).To(Equal("old.txt"))
					Expect(v.Newname).To(Equal("new.txt"))
					Expect(v.Path()).To(Equal("old.txt"))
					renameSeen = true
					return true
				case filter.CreateOp:
					// Rename also checks if the new name is allowed
					Expect(v.Name).To(Equal("new.txt"))
					createSeen = true
					return true
				default:
					return true
				}
			})

			err := filtered.Rename("old.txt", "new.txt")
			Expect(err).NotTo(HaveOccurred())
			Expect(renameSeen).To(BeTrue())
			Expect(createSeen).To(BeTrue())
		})

		It("should handle StatOp with type switch", func() {
			base := afero.NewMemMapFs()
			Expect(afero.WriteFile(base, "test.txt", []byte("test"), os.ModePerm)).To(Succeed())

			filtered := filter.NewFs(base, func(op filter.Operation) bool {
				switch v := op.(type) {
				case filter.StatOp:
					Expect(v.Name).To(Equal("test.txt"))
					Expect(v.Path()).To(Equal("test.txt"))
					return true
				default:
					return true
				}
			})

			_, err := filtered.Stat("test.txt")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should handle ReaddirOp with type switch", func() {
			base := afero.NewMemMapFs()
			Expect(base.Mkdir("dir", os.ModePerm)).To(Succeed())
			Expect(afero.WriteFile(base, "dir/file.txt", []byte("test"), os.ModePerm)).To(Succeed())

			filtered := filter.NewFs(base, func(op filter.Operation) bool {
				switch v := op.(type) {
				case filter.ReaddirOp:
					Expect(v.Name).To(ContainSubstring("dir"))
					Expect(v.Count).To(Equal(-1))
					Expect(v.Path()).To(ContainSubstring("dir"))
					return true
				default:
					return true
				}
			})

			file, err := filtered.Open("dir")
			Expect(err).NotTo(HaveOccurred())
			defer file.Close()

			_, err = file.Readdir(-1)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("PathFilter and PathPredicate", func() {
		It("should use PathFilter for backward compatibility", func() {
			base := afero.NewMemMapFs()
			Expect(afero.WriteFile(base, "test.txt", []byte("test"), os.ModePerm)).To(Succeed())

			pathFilter := filter.PathFilter(func(path string) error {
				if path == "test.txt" {
					return nil
				}
				return fs.ErrPermission
			})

			filtered := filter.FromFilter(base, pathFilter)
			_, err := filtered.Open("test.txt")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should block operations with PathPredicate", func() {
			base := afero.NewMemMapFs()
			Expect(afero.WriteFile(base, "test.txt", []byte("test"), os.ModePerm)).To(Succeed())

			pathPred := filter.PathPredicate(func(path string) bool {
				return path != "test.txt"
			})

			filtered := filter.NewFs(base, pathPred)
			_, err := filtered.Open("test.txt")
			Expect(err).To(HaveOccurred())
		})
	})
})
