package writer_test

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/unmango/aferox/writer"
)

var _ = Describe("Fs", func() {
	It("should write to provided writer", func() {
		buf := &bytes.Buffer{}
		fs := writer.NewFs(buf)

		file, err := fs.Open("doesn't matter")

		Expect(err).NotTo(HaveOccurred())
		_, err = file.WriteString("blahblahblah")
		Expect(err).NotTo(HaveOccurred())
		Expect(buf.String()).To(Equal("blahblahblah"))
	})
})
