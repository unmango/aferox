package asset_test

import (
	"testing"

	"github.com/google/go-github/v69/github"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unmango/aferox/github/internal"
)

var client *github.Client

func TestRelease(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Asset Suite")
}

var _ = BeforeSuite(func() {
	client = internal.DefaultClient()
})
