package user_test

import (
	"os"

	"github.com/google/go-github/v74/github"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/unmango/aferox/github/user"
)

var _ = Describe("Fs", func() {
	It("should open user", func() {
		client := github.NewClient(nil)
		fs := user.NewFs(client)

		user, err := fs.Open("UnstoppableMango")

		Expect(err).NotTo(HaveOccurred())
		Expect(user.Name()).To(Equal("UnstoppableMango"))
	})

	It("should open user file", func() {
		client := github.NewClient(nil)
		fs := user.NewFs(client)

		user, err := fs.OpenFile("UnstoppableMango", 69, os.ModePerm)

		Expect(err).NotTo(HaveOccurred())
		Expect(user.Name()).To(Equal("UnstoppableMango"))
	})

	It("should stat user", func() {
		client := github.NewClient(nil)
		fs := user.NewFs(client)

		user, err := fs.Stat("UnstoppableMango")

		Expect(err).NotTo(HaveOccurred())
		Expect(user.Name()).To(Equal("UnstoppableMango"))
	})

	It("should be readonly", func() {
		fs := user.NewFs(client)

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
})
