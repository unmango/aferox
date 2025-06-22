package mapped_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/unmango/aferox/mapped"
)

var _ = Describe("Path", func() {
	FDescribe("CutPrefix", func() {
		DescribeTable("should cut matching prefixes",
			func(s, prefix, expected string) {
				actual, found := mapped.CutPrefix(s, prefix)

				Expect(found).To(BeTrueBecause("The prefix was found"))
				Expect(actual).To(Equal(expected))
			},
			Entry(nil, "test/test.txt", "test", "/test.txt"),
			Entry(nil, "/test/test.txt", "/test", "/test.txt"),
			Entry(nil, "test/test.txt", "/test", "/test.txt"),
		)
	})
})
