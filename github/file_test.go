package github_test

import (
	"io"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	ghx "github.com/unmango/aferox/github"
)

var _ = Describe("File", func() {
	It("should list repositories", func() {
		fs := ghx.NewFs(client)
		file, err := fs.Open("users/UnstoppableMango")
		Expect(err).NotTo(HaveOccurred())

		names, err := file.Readdirnames(69)

		Expect(err).NotTo(HaveOccurred())
		Expect(names).To(ContainElement("advent-of-code"))
	})

	It("should read json", func() {
		fs := ghx.NewFs(client)
		file, err := fs.Open("users/UnstoppableMango")
		Expect(err).NotTo(HaveOccurred())

		data, err := io.ReadAll(file)

		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).To(And(
			ContainSubstring("login"),
			ContainSubstring("UnstoppableMango"),
		))
		Expect(file.Close()).To(Succeed())
	})

	It("should Open user", func() {
		fs := ghx.NewFs(client)

		file, err := fs.Open("users/UnstoppableMango")

		Expect(err).NotTo(HaveOccurred())
		Expect(file.Name()).To(Equal("users/UnstoppableMango"))
	})

	It("should be readonly", func() {
		fs := ghx.NewFs(client)
		file, err := fs.Open("users/UnstoppableMango")
		Expect(err).NotTo(HaveOccurred())

		_, err = file.Write([]byte{})
		Expect(err).To(MatchError("read-only file system"))
		_, err = file.WriteAt([]byte{}, 69)
		Expect(err).To(MatchError("read-only file system"))
		_, err = file.WriteString("doesn't matter")
		Expect(err).To(MatchError("read-only file system"))
	})

	It("should read archive as directory", Label("E2E"), func() {
		fs := ghx.NewFs(client)
		file, err := fs.Open("github.com/UnstoppableMango/tdl/releases/tag/v0.0.29/tdl-linux-amd64.tar.gz")
		Expect(err).NotTo(HaveOccurred())

		stat, err := file.Stat()
		Expect(err).NotTo(HaveOccurred())
		Expect(stat.IsDir()).To(BeTrueBecause("treat archives as directories"))
		names, err := file.Readdirnames(3)
		Expect(err).NotTo(HaveOccurred())
		Expect(names).To(ConsistOf("uml2ts", "uml2go", "uml2pcl"))
	})
})
