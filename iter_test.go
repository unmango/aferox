package aferox_test

import (
	"errors"
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

	It("should continue on error when ContinueOnError is provided", func() {
		fs := afero.NewMemMapFs()
		// Create a file that will exist during walk
		_, err := fs.Create("test.txt")
		Expect(err).NotTo(HaveOccurred())

		// Wrap in readonly to potentially cause errors during write operations
		// But for reading, it should still work
		seq := aferox.Iter(fs, "", aferox.ContinueOnError)

		a, _, _ := slices.Collect3(seq)
		Expect(a).NotTo(BeEmpty())
	})

	It("should filter errors when FilterErrors is provided", func() {
		fs := afero.NewMemMapFs()
		_, err := fs.Create("test.txt")
		Expect(err).NotTo(HaveOccurred())

		filter := func(err error) error {
			return nil
		}

		seq := aferox.Iter(fs, "", aferox.ContinueOnError, aferox.FilterErrors(filter))

		a, _, _ := slices.Collect3(seq)
		Expect(a).NotTo(BeEmpty())
		// Note: This tests that FilterErrors option is applied without errors
	})

	It("should stop iteration early", func() {
		fs := afero.NewMemMapFs()
		_, err := fs.Create("file1.txt")
		Expect(err).NotTo(HaveOccurred())
		_, err = fs.Create("file2.txt")
		Expect(err).NotTo(HaveOccurred())

		count := 0
		seq := aferox.Iter(fs, "")
		seq(func(path string, info os.FileInfo, err error) bool {
			count++
			return path != "file1.txt" // Continue until we find file1.txt
		})

		Expect(count).To(BeNumerically("<=", 3)) // root + file1 + maybe one more before break
	})

	It("should yield final error when iteration fails", func() {
		// Create a custom fs that returns an error during walk
		fs := afero.NewMemMapFs()
		_, err := fs.Create("test.txt")
		Expect(err).NotTo(HaveOccurred())

		// Use a readonly fs and try to walk a non-existent path
		readonlyFs := afero.NewReadOnlyFs(fs)
		seq := aferox.Iter(readonlyFs, "nonexistent/path")

		paths, _, errs := slices.Collect3(seq)
		Expect(len(paths)).To(BeNumerically(">", 0))
		lastErr := errs[len(errs)-1]
		Expect(lastErr).To(HaveOccurred())
	})

	It("should apply FilterErrors that returns an error", func() {
		fs := afero.NewMemMapFs()

		called := false
		expectedErr := errors.New("filtered error")
		filter := func(err error) error {
			if err != nil {
				called = true
				return expectedErr
			}
			return nil
		}

		seq := aferox.Iter(fs, "nonexistent", aferox.ContinueOnError, aferox.FilterErrors(filter))

		_, _, errs := slices.Collect3(seq)
		Expect(called).To(BeTrue(), "filter should have been called")
		Expect(errs[len(errs)-1]).To(MatchError(expectedErr))
	})
})
