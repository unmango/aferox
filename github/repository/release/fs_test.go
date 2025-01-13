package release_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/unmango/aferox/github/repository/release"
)

var _ = Describe("Fs", func() {
	It("should stat", func() {
		fs := release.NewFs(client, "UnstoppableMango", "tdl")

		r, err := fs.Stat("releases/tag/v0.0.29")

		Expect(err).NotTo(HaveOccurred())
		Expect(r.Name()).To(Equal("v0.0.29"))
	})

	It("should stat asset", func() {
		fs := release.NewFs(client, "UnstoppableMango", "tdl")

		r, err := fs.Stat("releases/tag/v0.0.29/tdl-linux-amd64.tar.gz")

		Expect(err).NotTo(HaveOccurred())
		Expect(r.Name()).To(Equal("tdl-linux-amd64.tar.gz"))
	})
})
