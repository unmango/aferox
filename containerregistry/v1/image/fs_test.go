package image_test

import (
	"archive/tar"
	"errors"
	"io"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/spf13/afero"
	"github.com/unmango/aferox/containerregistry/v1/image"
)

// https://github.com/google/go-containerregistry/blob/main/pkg/crane/filemap_test.go
var _ = DescribeTableSubtree("Fs",
	func(memfs map[string][]byte, digest string) {
		var fs afero.Fs

		BeforeEach(func() {
			fs = afero.NewMemMapFs()
			for path, content := range memfs {
				err := afero.WriteFile(fs, path, content, 0644)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		Describe("ToFs", func() {
			It("should match contents", func() {
				l, err := crane.Image(memfs)
				Expect(err).NotTo(HaveOccurred())

				fs, err = image.ToFs(l)

				Expect(err).NotTo(HaveOccurred())
				saw := map[string]struct{}{}
				_ = afero.Walk(fs, "/", func(path string, info os.FileInfo, err error) error {
					if path == "/" {
						return nil
					}

					Expect(err).NotTo(HaveOccurred())
					saw[path] = struct{}{}
					want, found := memfs[path]
					Expect(found).To(BeTrueBecause("found %q, not in original map", path))

					got, err := afero.ReadFile(fs, path)
					Expect(err).NotTo(HaveOccurred())
					Expect(string(got)).To(Equal(string(want)))

					return nil
				})
			})
		})

		Describe("FromFs", func() {
			It("should match digest", func() {
				img, err := image.FromFs(fs)
				Expect(err).NotTo(HaveOccurred())

				d, err := img.Digest()
				Expect(err).NotTo(HaveOccurred())
				Expect(d.String()).To(Equal(digest))
			})

			It("should match contents", func() {
				img, err := image.FromFs(fs)
				Expect(err).NotTo(HaveOccurred())

				rc := mutate.Extract(img)
				defer rc.Close()

				tr := tar.NewReader(rc)
				saw := map[string]struct{}{}
				for {
					th, err := tr.Next()
					if errors.Is(err, io.EOF) {
						break
					}

					Expect(err).NotTo(HaveOccurred())
					saw[th.Name] = struct{}{}
					want, found := memfs[th.Name]
					Expect(found).To(BeTrueBecause("found %q, not in original map", th.Name))

					got, err := io.ReadAll(tr)
					Expect(err).NotTo(HaveOccurred())
					Expect(string(got)).To(Equal(string(want)))
				}
			})

			It("should be reproducible", func() {
				i1, err := image.FromFs(fs)
				Expect(err).NotTo(HaveOccurred())

				i2, err := image.FromFs(fs)
				Expect(err).NotTo(HaveOccurred())

				d1, err := i1.Digest()
				Expect(err).NotTo(HaveOccurred())
				d2, err := i2.Digest()
				Expect(err).NotTo(HaveOccurred())
				Expect(d2.String()).To(Equal(d1.String()))
			})
		})
	},
	Entry("Empty contents",
		map[string][]byte{},
		"sha256:98132f58b523c391a5788997327cac95e114e3a6609d01163189774510705399",
	),
	Entry("One file",
		map[string][]byte{
			"/test": []byte("testy"),
		},
		"sha256:d905c03ac635172a96c12b8af6c90cfd028e3edaa3114b31a9e196ab38c16963",
	),
	Entry("Two files",
		map[string][]byte{
			"/test": []byte("testy"),
			"/bar":  []byte("not useful"),
		},
		"sha256:20e7e4800e5eb167f170970936c08d9e1bcbe91372420eeb6ab8d1a07752c3a3",
	),
	Entry("Many files",
		map[string][]byte{
			"/1": []byte("1"),
			"/2": []byte("2"),
			"/3": []byte("3"),
			"/4": []byte("4"),
			"/5": []byte("5"),
			"/6": []byte("6"),
			"/7": []byte("7"),
			"/8": []byte("8"),
			"/9": []byte("9"),
		},
		"sha256:dfca2803510c8e3b83a3151f7c035c60cfa2a8a52465b802e18b85014de361f1",
	),
)
