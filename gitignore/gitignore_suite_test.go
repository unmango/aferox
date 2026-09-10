package gitignore_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGitIgnore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GitIgnore Suite")
}
