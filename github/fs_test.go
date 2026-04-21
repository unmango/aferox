package github_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	ghx "github.com/unmango/aferox/github"
)

var _ = Describe("Fs", func() {
	// User
	It("should open user", func() {
		fs := ghx.NewFs(client)

		file, err := fs.Open("users/UnstoppableMango")

		Expect(err).NotTo(HaveOccurred())
		Expect(file.Name()).To(Equal("users/UnstoppableMango"))
	})

	It("should open user file", func() {
		fs := ghx.NewFs(client)

		file, err := fs.OpenFile("users/UnstoppableMango", 69, os.ModePerm)

		Expect(err).NotTo(HaveOccurred())
		Expect(file.Name()).To(Equal("users/UnstoppableMango"))
	})

	It("should stat user", func() {
		fs := ghx.NewFs(client)

		info, err := fs.Stat("users/UnstoppableMango")

		Expect(err).NotTo(HaveOccurred())
		Expect(info.Name()).To(Equal("UnstoppableMango"))
	})

	It("should be readonly", func() {
		fs := ghx.NewFs(client)

		_, err := fs.Create("doesn't matter")
		Expect(err).To(MatchError("operation not permitted"))
		err = fs.Chmod("doesn't matter", os.ModeSetgid)
		Expect(err).To(MatchError("operation not permitted"))
		err = fs.Chown("doesn't matter", 420, 69)
		Expect(err).To(MatchError("operation not permitted"))
		err = fs.Mkdir("doesn't matter", os.ModeDir)
		Expect(err).To(MatchError("operation not permitted"))
		err = fs.MkdirAll("doesn't matter", os.ModeDir)
		Expect(err).To(MatchError("operation not permitted"))
		err = fs.Remove("doesn't matter")
		Expect(err).To(MatchError("operation not permitted"))
		err = fs.RemoveAll("doesn't matter")
		Expect(err).To(MatchError("operation not permitted"))
		err = fs.Rename("doesn't matter", "still doesn't matter")
		Expect(err).To(MatchError("operation not permitted"))
	})

	// Repository
	It("should open repo", func() {
		fs := ghx.NewFs(client)

		repo, err := fs.Open("repos/UnstoppableMango/advent-of-code")

		Expect(err).NotTo(HaveOccurred())
		Expect(repo.Name()).To(Equal("repos/UnstoppableMango/advent-of-code"))
	})

	It("should stat release", func() {
		fs := ghx.NewFs(client)

		info, err := fs.Stat("repos/UnstoppableMango/tdl/releases/tags/v0.0.29")

		Expect(err).NotTo(HaveOccurred())
		Expect(info.Name()).To(Equal("v0.0.29"))
	})

	// Content
	It("should stat file", func() {
		fs := ghx.NewFs(client)

		info, err := fs.Stat("github.com/UnstoppableMango/tdl/tree/main/Makefile")

		Expect(err).NotTo(HaveOccurred())
		Expect(info.Name()).To(Equal("Makefile"))
	})

	It("should open file", func() {
		fs := ghx.NewFs(client)

		file, err := fs.Open("github.com/UnstoppableMango/tdl/tree/main/Makefile")

		Expect(err).NotTo(HaveOccurred())
		Expect(file.Name()).To(Equal("github.com/UnstoppableMango/tdl/tree/main/Makefile"))
	})

	It("should open directory", func() {
		fs := ghx.NewFs(client)

		file, err := fs.Open("github.com/UnstoppableMango/tdl/tree/main/cmd")

		Expect(err).NotTo(HaveOccurred())
		Expect(file.Name()).To(Equal("github.com/UnstoppableMango/tdl/tree/main/cmd"))
		Expect(file.Readdirnames(3)).To(ConsistOf("ux", "uml2uml"))
	})

	// Asset
	It("should stat asset", func() {
		fs := ghx.NewFs(client)

		info, err := fs.Stat("github.com/UnstoppableMango/tdl/releases/tag/v0.0.29/tdl-linux-amd64.tar.gz")

		Expect(err).NotTo(HaveOccurred())
		Expect(info.IsDir()).To(BeTrueBecause("support reading archives"))
		Expect(info.Name()).To(Equal("tdl-linux-amd64.tar.gz"))
	})

	It("should download an asset", Label("E2E"), func() {
		fs := ghx.NewFs(client)

		file, err := fs.Open("github.com/UnstoppableMango/tdl/releases/tag/v0.0.29/tdl-linux-amd64.tar.gz")

		Expect(err).NotTo(HaveOccurred())
		Expect(file.Name()).To(Equal("github.com/UnstoppableMango/tdl/releases/tag/v0.0.29/tdl-linux-amd64.tar.gz"))
	})
})
