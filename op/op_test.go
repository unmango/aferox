package op_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/unmango/aferox/op"
)

var _ = Describe("Op", func() {
	Describe("Operation interface", func() {
		It("should implement Operation for Chmod", func() {
			var operation op.Operation = op.Chmod{
				Name: "test.txt",
				Mode: 0755,
			}

			Expect(operation.Path()).To(Equal("test.txt"))
		})

		It("should implement Operation for Chown", func() {
			var operation op.Operation = op.Chown{
				Name: "test.txt",
				UID:  1000,
				GID:  1000,
			}

			Expect(operation.Path()).To(Equal("test.txt"))
		})

		It("should implement Operation for Chtimes", func() {
			now := time.Now()
			var operation op.Operation = op.Chtimes{
				Name:  "test.txt",
				Atime: now,
				Mtime: now,
			}

			Expect(operation.Path()).To(Equal("test.txt"))
		})

		It("should implement Operation for Create", func() {
			var operation op.Operation = op.Create{Name: "test.txt"}

			Expect(operation.Path()).To(Equal("test.txt"))
		})

		It("should implement Operation for Mkdir", func() {
			var operation op.Operation = op.Mkdir{
				Name: "testdir",
				Perm: 0755,
			}

			Expect(operation.Path()).To(Equal("testdir"))
		})

		It("should implement Operation for MkdirAll", func() {
			var operation op.Operation = op.MkdirAll{
				PathName: "path/to/dir",
				Perm:     0755,
			}

			Expect(operation.Path()).To(Equal("path/to/dir"))
		})

		It("should implement Operation for Open", func() {
			var operation op.Operation = op.Open{Name: "test.txt"}

			Expect(operation.Path()).To(Equal("test.txt"))
		})

		It("should implement Operation for OpenFile", func() {
			var operation op.Operation = op.OpenFile{
				Name: "test.txt",
				Flag: 0,
				Perm: 0644,
			}

			Expect(operation.Path()).To(Equal("test.txt"))
		})

		It("should implement Operation for Remove", func() {
			var operation op.Operation = op.Remove{Name: "test.txt"}

			Expect(operation.Path()).To(Equal("test.txt"))
		})

		It("should implement Operation for RemoveAll", func() {
			var operation op.Operation = op.RemoveAll{PathName: "testdir"}

			Expect(operation.Path()).To(Equal("testdir"))
		})

		It("should implement Operation for Rename", func() {
			var operation op.Operation = op.Rename{
				Oldname: "old.txt",
				Newname: "new.txt",
			}

			Expect(operation.Path()).To(Equal("old.txt"))
		})

		It("should implement Operation for Stat", func() {
			var operation op.Operation = op.Stat{Name: "test.txt"}

			Expect(operation.Path()).To(Equal("test.txt"))
		})

		It("should implement Operation for Readdir", func() {
			var operation op.Operation = op.Readdir{
				Name:  "dir",
				Count: -1,
			}

			Expect(operation.Path()).To(Equal("dir"))
		})

		It("should implement Operation for Readdirnames", func() {
			var operation op.Operation = op.Readdirnames{
				Name:  "dir",
				Count: -1,
			}

			Expect(operation.Path()).To(Equal("dir"))
		})
	})
})
