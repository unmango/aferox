package github_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/google/go-github/v84/github"
	"github.com/unmango/aferox/github/internal"
)

var client *github.Client

func TestGithub(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Github Suite")
}

var _ = BeforeSuite(func() {
	client = internal.DefaultClient()
})
