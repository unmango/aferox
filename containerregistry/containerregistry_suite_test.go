package containerregistry_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestContainerRegistry(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ContainerRegistry Suite")
}
