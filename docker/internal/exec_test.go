package internal_test

import (
	"bytes"
	"context"

	"github.com/docker/docker/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"

	"github.com/unmango/aferox/docker/internal"
)

type logger struct{}

func (l logger) Printf(format string, v ...interface{}) {
	GinkgoWriter.Printf(format+"\n", v)
}

var _ = Describe("Exec", Label("E2E"), func() {
	var (
		ctr    testcontainers.Container
		docker client.APIClient
	)

	BeforeEach(func(ctx context.Context) {
		req := testcontainers.ContainerRequest{
			Image: "ubuntu",
			Cmd:   []string{"sleep", "infinity"},
		}

		var err error
		ctr, err = testcontainers.GenericContainer(ctx,
			testcontainers.GenericContainerRequest{
				ContainerRequest: req,
				Started:          true,
				Logger:           logger{},
			},
		)
		Expect(err).NotTo(HaveOccurred())

		docker, err = client.NewClientWithOpts(
			client.WithAPIVersionNegotiation(),
		)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := testcontainers.TerminateContainer(ctr)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should work", func(ctx context.Context) {
		buf := &bytes.Buffer{}
		err := internal.Exec(ctx, docker, ctr.GetContainerID(), internal.ExecOptions{
			Cmd:    []string{"echo", "testing"},
			Stdout: buf,
		})

		Expect(err).NotTo(HaveOccurred())
		Expect(buf.String()).To(Equal("testing\n"))
	})
})
