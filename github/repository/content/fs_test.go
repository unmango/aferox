package content_test

import (
	"io"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/unmango/aferox/github/repository/content"
)

var _ = Describe("Fs", func() {
	It("should stat file", func() {
		fs := content.NewFs(client, "UnstoppableMango", "tdl", "main")

		stat, err := fs.Stat("Makefile")

		Expect(err).NotTo(HaveOccurred())
		Expect(stat.Name()).To(Equal("Makefile"))
	})

	It("should open file", func() {
		fs := content.NewFs(client, "UnstoppableMango", "tdl", "main")

		file, err := fs.Open("Makefile")

		Expect(err).NotTo(HaveOccurred())
		Expect(file.Name()).To(Equal("Makefile"))
		data, err := io.ReadAll(file)
		Expect(data).NotTo(BeEmpty())
	})

	It("should open directory", func() {
		fs := content.NewFs(client, "UnstoppableMango", "tdl", "main")

		file, err := fs.Open("cmd")

		Expect(err).NotTo(HaveOccurred())
		Expect(file.Name()).To(Equal("cmd"))
		Expect(file.Readdirnames(3)).To(
			ConsistOf("ux", "uml2uml"),
		)
	})
})
