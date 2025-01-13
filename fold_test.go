package aferox_test

import (
	"io/fs"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"

	aferox "github.com/unmango/aferox"
)

var _ = Describe("Fold", func() {
	It("should work", func() {
		fsys := afero.NewMemMapFs()
		_, err := fsys.Create("test.txt")
		Expect(err).NotTo(HaveOccurred())
		var count int

		res, err := aferox.Fold(fsys, "",
			func(path string, info fs.FileInfo, acc int, err error) (int, error) {
				return acc + 1, err
			},
			count,
		)

		Expect(err).NotTo(HaveOccurred())
		// Includes the "." path
		Expect(res).To(Equal(2))
	})
})
