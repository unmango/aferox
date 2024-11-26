package docker_test

import (
	"context"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	. "github.com/unmango/go/testing/matcher"

	"github.com/unmango/go/fs/docker"
)

var _ = Describe("Buffered", func() {
	It("should work", func(ctx context.Context) {
		fs, sync := docker.NewBufferedFs(testclient, ctr.GetContainerID())
		err := afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		Expect(fs).NotTo(ContainFile("test.txt"))

		Expect(sync(ctx)).To(Succeed())

		Expect(fs).To(ContainFileWithBytes("test.txt", []byte("testing")))
	})
})
