package aferox_test

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/unmango/aferox"
)

var _ = Describe("Fs", func() {
	Describe("NewWriter", func() {
		It("should create a Writer Fs", func() {
			var buf bytes.Buffer

			fs := aferox.NewWriter(&buf)

			Expect(fs).NotTo(BeNil())
			Expect(fs.Name()).To(Equal("io.Writer"))
		})

		It("should return a usable Fs", func() {
			var buf bytes.Buffer
			fs := aferox.NewWriter(&buf)

			file, err := fs.Open("test.txt")

			Expect(err).NotTo(HaveOccurred())
			Expect(file).NotTo(BeNil())
		})
	})
})
