package repository_test

import (
	"testing"

	"github.com/google/go-github/v82/github"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unmango/aferox/github/internal"
)

var client *github.Client

func TestRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repository Suite")
}

var _ = BeforeSuite(func() {
	client = internal.DefaultClient()
})
