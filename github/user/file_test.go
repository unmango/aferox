package user_test

import (
	"io"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/unmango/aferox/github/user"
)

var _ = Describe("File", func() {
	It("should list repositories", func() {
		fs := user.NewFs(client)
		file, err := fs.Open("UnstoppableMango")
		Expect(err).NotTo(HaveOccurred())

		names, err := file.Readdirnames(69)

		Expect(err).NotTo(HaveOccurred())
		Expect(names).To(ContainElement("advent-of-code"))
	})

	It("should read json", func() {
		fs := user.NewFs(client)
		file, err := fs.Open("UnstoppableMango")
		Expect(err).NotTo(HaveOccurred())

		data, err := io.ReadAll(file)

		Expect(err).NotTo(HaveOccurred())
		Expect(data).To(And(
			ContainSubstring("login"),
			ContainSubstring("UnstoppableMango"),
		))
		Expect(file.Close()).To(Succeed())
	})

	It("should Open user", func() {
		fs := user.NewFs(client)

		file, err := fs.Open("UnstoppableMango")

		Expect(err).NotTo(HaveOccurred())
		Expect(file.Name()).To(Equal("UnstoppableMango"))
	})

	It("should be readonly", func() {
		fs := user.NewFs(client)
		file, err := fs.Open("UnstoppableMango")
		Expect(err).NotTo(HaveOccurred())

		_, err = file.Write([]byte{})
		Expect(err).To(MatchError("read-only file system"))
		_, err = file.WriteAt([]byte{}, 69)
		Expect(err).To(MatchError("read-only file system"))
		_, err = file.WriteString("doesn't matter")
		Expect(err).To(MatchError("read-only file system"))
	})
})
