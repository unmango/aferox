package docker_test

import (
	"context"
	"io"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/unmango/aferox/docker"
)

var _ = Describe("Fs", Label("E2E"), func() {
	It("should list directories", func(ctx context.Context) {
		fs := docker.NewFs(testclient, ctr.GetContainerID())
		dir, err := fs.Open(ctx, "/")
		Expect(err).NotTo(HaveOccurred())

		infos, err := dir.Readdir(69)

		Expect(err).NotTo(HaveOccurred())
		names := make([]string, len(infos))
		for i, f := range infos {
			names[i] = f.Name()
		}
		Expect(names).To(ContainElements("root", "var", "bin"))
	})

	It("should read file contents", func(ctx context.Context) {
		fs := docker.NewFs(testclient, ctr.GetContainerID())
		file, err := fs.Create(ctx, "test-read.txt")
		Expect(err).NotTo(HaveOccurred())
		_, err = io.WriteString(file, "bleh")
		Expect(err).NotTo(HaveOccurred())
		Expect(file.Close()).To(Succeed())

		file, err = fs.Open(ctx, "test-read.txt")

		Expect(err).NotTo(HaveOccurred())
		data, err := io.ReadAll(file)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).To(Equal("bleh"))
	})

	Describe("Create", func() {
		It("should create a file", func(ctx context.Context) {
			fsys := docker.NewFs(testclient, ctr.GetContainerID())

			file, err := fsys.Create(ctx, "test.txt")

			Expect(err).NotTo(HaveOccurred())
			Expect(file).NotTo(BeNil())
		})

		It("should create a writable file", func(ctx context.Context) {
			fsys := docker.NewFs(testclient, ctr.GetContainerID())

			file, err := fsys.Create(ctx, "writable.txt")

			Expect(err).NotTo(HaveOccurred())
			_, err = file.WriteString("blahblahblah")
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
