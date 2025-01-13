package asset_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unmango/aferox/github/repository/release/asset"
)

var _ = Describe("File", func() {
	It("should be readonly", func() {
		fs := asset.NewFs(client, "UnstoppableMango", "tdl", "v0.0.29")
		file, err := fs.Open("tdl-linux-amd64.tar.gz")
		Expect(err).NotTo(HaveOccurred())

		_, err = file.Write([]byte{})
		Expect(err).To(MatchError("read-only file system"))
		_, err = file.WriteAt([]byte{}, 69)
		Expect(err).To(MatchError("read-only file system"))
		_, err = file.WriteString("doesn't matter")
		Expect(err).To(MatchError("read-only file system"))
	})
})
