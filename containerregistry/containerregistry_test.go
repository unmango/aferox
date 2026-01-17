package containerregistry_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/spf13/afero"
	"github.com/unmango/aferox/containerregistry"
)

var _ = Describe("Exports", func() {
	var fs afero.Fs

	BeforeEach(func() {
		fs = afero.NewMemMapFs()
		err := afero.WriteFile(fs, "/test", []byte("test content"), 0644)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("ImageFromFs", func() {
		It("should create an image from filesystem", func() {
			img, err := containerregistry.ImageFromFs(fs)
			Expect(err).NotTo(HaveOccurred())
			Expect(img).NotTo(BeNil())

			// Verify the image has layers
			layers, err := img.Layers()
			Expect(err).NotTo(HaveOccurred())
			Expect(layers).To(HaveLen(1))
		})
	})

	Describe("ImageToFs", func() {
		It("should extract an image to filesystem", func() {
			img, err := containerregistry.ImageFromFs(fs)
			Expect(err).NotTo(HaveOccurred())

			extractedFs, err := containerregistry.ImageToFs(img)
			Expect(err).NotTo(HaveOccurred())
			Expect(extractedFs).NotTo(BeNil())

			// Verify the content was extracted
			content, err := afero.ReadFile(extractedFs, "/test")
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(Equal("test content"))
		})
	})

	Describe("LayerFromFs", func() {
		It("should create a layer from filesystem", func() {
			layer, err := containerregistry.LayerFromFs(fs)
			Expect(err).NotTo(HaveOccurred())
			Expect(layer).NotTo(BeNil())

			// Verify the layer has content
			size, err := layer.Size()
			Expect(err).NotTo(HaveOccurred())
			Expect(size).To(BeNumerically(">", 0))
		})
	})

	Describe("LayerToFs", func() {
		It("should extract a layer to filesystem", func() {
			layer, err := containerregistry.LayerFromFs(fs)
			Expect(err).NotTo(HaveOccurred())

			extractedFs, err := containerregistry.LayerToFs(layer)
			Expect(err).NotTo(HaveOccurred())
			Expect(extractedFs).NotTo(BeNil())

			// Verify the content was extracted
			content, err := afero.ReadFile(extractedFs, "/test")
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(Equal("test content"))
		})
	})
})
