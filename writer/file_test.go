package writer_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/unmango/aferox/writer"
)

type stub struct {
	closed bool
}

func (s *stub) Write(p []byte) (int, error) {
	return 0, nil
}

func (s *stub) Close() error {
	s.closed = true
	return nil
}

var _ = Describe("File", func() {
	It("should not close the parent writer", func() {
		w := &stub{}
		fs := writer.NewFs(w)
		file, err := fs.Open("doesn't matter")
		Expect(err).NotTo(HaveOccurred())

		err = file.Close()

		Expect(err).NotTo(HaveOccurred())
		Expect(w.closed).To(BeFalseBecause("the writer is not closed"))
	})
})
