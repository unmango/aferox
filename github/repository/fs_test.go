package repository_test

import (
	"github.com/google/go-github/v82/github"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/unmango/aferox/github/repository"
)

var _ = Describe("Fs", func() {
	It("should open repo", func() {
		client := github.NewClient(nil)
		fs := repository.NewFs(client, "UnstoppableMango")

		repo, err := fs.Open("advent-of-code")

		Expect(err).NotTo(HaveOccurred())
		Expect(repo.Name()).To(Equal("advent-of-code"))
	})

	It("should stat release", func() {
		client := github.NewClient(nil)
		fs := repository.NewFs(client, "UnstoppableMango")

		release, err := fs.Stat("tdl/releases/tag/v0.0.29")

		Expect(err).NotTo(HaveOccurred())
		Expect(release.Name()).To(Equal("v0.0.29"))
	})
})
