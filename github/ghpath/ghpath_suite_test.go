package ghpath_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGhpath(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ghpath Suite")
}
