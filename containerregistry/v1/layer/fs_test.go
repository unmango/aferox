package layer_test

import (
	"archive/tar"
	"errors"
	"io"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/spf13/afero"
	"github.com/unmango/aferox/containerregistry/v1/layer"
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
				l, err := crane.Layer(memfs)
				Expect(err).NotTo(HaveOccurred())

				fs, err = layer.ToFs(l)

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
				l, err := layer.FromFs(fs)
				Expect(err).NotTo(HaveOccurred())

				d, err := l.Digest()
				Expect(err).NotTo(HaveOccurred())
				Expect(d.String()).To(Equal(digest))
			})

			It("should match contents", func() {
				l, err := layer.FromFs(fs)
				Expect(err).NotTo(HaveOccurred())

				rc, err := l.Uncompressed()
				Expect(err).NotTo(HaveOccurred())
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
				l1, err := layer.FromFs(fs)
				Expect(err).NotTo(HaveOccurred())

				l2, err := layer.FromFs(fs)
				Expect(err).NotTo(HaveOccurred())

				d1, err := l1.Digest()
				Expect(err).NotTo(HaveOccurred())
				d2, err := l2.Digest()
				Expect(err).NotTo(HaveOccurred())
				Expect(d2.String()).To(Equal(d1.String()))
			})
		})
	},
	Entry("Empty contents",
		map[string][]byte{},
		"sha256:89732bc7504122601f40269fc9ddfb70982e633ea9caf641ae45736f2846b004",
	),
	Entry("One file",
		map[string][]byte{
			"/test": []byte("testy"),
		},
		"sha256:ec3ff19f471b99a76fb1c339c1dfdaa944a4fba25be6bcdc99fe7e772103079e",
	),
	Entry("Two files",
		map[string][]byte{
			"/test":    []byte("testy"),
			"/testalt": []byte("footesty"),
		},
		"sha256:a48bcb7be3ab3ec608ee56eb80901224e19e31dc096cc06a8fd3a8dae1aa8947",
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
		"sha256:1e637602abbcab2dcedcc24e0b7c19763454a47261f1658b57569530b369ccb9",
	),
)
