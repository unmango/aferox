package asset_test

import (
	"io"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/unmango/aferox/github/repository/release/asset"
)

var _ = Describe("Fs", func() {
	It("should stat an asset", func() {
		r := asset.NewFs(client, "UnstoppableMango", "tdl", "v0.0.29")

		info, err := r.Stat("tdl-linux-amd64.tar.gz")

		Expect(err).NotTo(HaveOccurred())
		Expect(info.IsDir()).To(BeTrueBecause("support reading archives"))
		Expect(info.Name()).To(Equal("tdl-linux-amd64.tar.gz"))
	})

	It("should download an asset", Label("E2E"), func() {
		r := asset.NewFs(client, "UnstoppableMango", "tdl", "v0.0.29")

		file, err := r.Open("tdl-linux-amd64.tar.gz")

		Expect(err).NotTo(HaveOccurred())
		Expect(file.Name()).To(Equal("tdl-linux-amd64.tar.gz"))
		data, err := io.ReadAll(file)
		Expect(err).NotTo(HaveOccurred())
		Expect(data).To(HaveLen(49_388_058))
	})

	It("should read an archive asset", Label("E2E"), func() {
		r := asset.NewFs(client, "UnstoppableMango", "tdl", "v0.0.29")

		file, err := r.Open("tdl-linux-amd64.tar.gz")

		Expect(err).NotTo(HaveOccurred())
		stat, err := file.Stat()
		Expect(err).NotTo(HaveOccurred())
		Expect(stat.IsDir()).To(BeTrueBecause("treat archives as directories"))
		infos, err := file.Readdirnames(3)
		Expect(err).NotTo(HaveOccurred())
		Expect(infos).To(ConsistOf("uml2ts", "uml2go", "uml2pcl"))
	})
})
