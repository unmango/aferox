package mapped_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMapped(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mapped Suite")
}
