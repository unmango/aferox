package layer_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLayer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "V1 Layer Suite")
}
