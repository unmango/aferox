package protofsv1alpha1_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGrpcV1alpha1(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gRPC v1alpha1 Suite")
}
