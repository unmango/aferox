package content_test

import (
	"io"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/unmango/aferox/github/repository/content"
)

var _ = Describe("Fs", func() {
	It("should stat file", func() {
		fs := content.NewFs(client, "unmango", "aferox", "main")

		stat, err := fs.Stat("Makefile")

		Expect(err).NotTo(HaveOccurred())
		Expect(stat.Name()).To(Equal("Makefile"))
	})

	It("should open file", func() {
		fs := content.NewFs(client, "unmango", "aferox", "main")

		file, err := fs.Open("Makefile")

		Expect(err).NotTo(HaveOccurred())
		Expect(file.Name()).To(Equal("Makefile"))
		data, err := io.ReadAll(file)
		Expect(data).NotTo(BeEmpty())
	})

	It("should open directory", func() {
		fs := content.NewFs(client, "unmango", "aferox", "main")

		file, err := fs.Open("github")

		Expect(err).NotTo(HaveOccurred())
		Expect(file.Name()).To(Equal("github"))
		Expect(file.Readdirnames(100)).To(
			// This is still terrible, but it's better than before
			ContainElements("ghpath", "repository", "user"),
		)
	})
})
