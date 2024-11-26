package docker_test

import (
	"context"
	"testing"

	"github.com/docker/docker/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
)

type logger struct{}

func (l logger) Printf(format string, v ...interface{}) {
	GinkgoWriter.Printf(format+"\n", v)
}

var (
	testclient client.ContainerAPIClient
	ctr        testcontainers.Container
)

func TestDocker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Docker Suite")
}

var _ = BeforeSuite(func(ctx context.Context) {
	var err error
	testclient, err = client.NewClientWithOpts(
		client.WithAPIVersionNegotiation(),
	)
	Expect(err).NotTo(HaveOccurred())

	req := testcontainers.ContainerRequest{
		Image: "ubuntu",
		Cmd:   []string{"sleep", "infinity"},
	}

	ctr, err = testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
			Logger:           logger{},
		},
	)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	err := testcontainers.TerminateContainer(ctr)
	Expect(err).NotTo(HaveOccurred())
})
