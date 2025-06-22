package mapped_test

import (
	"io"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/spf13/afero"
	"github.com/unmango/aferox/mapped"
)

var _ = Describe("Fs", func() {
	It("should open a file in the nested fs", func() {
		testFs := afero.NewMemMapFs()
		err := afero.WriteFile(testFs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		fs := mapped.NewFs(map[string]afero.Fs{
			"test": testFs,
		})

		f, err := fs.Open("test/test.txt")

		Expect(err).NotTo(HaveOccurred())
		data, err := io.ReadAll(f)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).To(Equal("testing"))
	})

	It("should open a file at the root of the nested fs", func() {
		testFs := afero.NewMemMapFs()
		err := afero.WriteFile(testFs, "/test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		fs := mapped.NewFs(map[string]afero.Fs{
			"test": testFs,
		})

		f, err := fs.Open("/test/test.txt")

		Expect(err).NotTo(HaveOccurred())
		data, err := io.ReadAll(f)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).To(Equal("testing"))
	})

	It("should open a file in an fs with multiple path segments", func() {
		testFs := afero.NewMemMapFs()
		err := afero.WriteFile(testFs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		fs := mapped.NewFs(map[string]afero.Fs{
			"test/with-segment": testFs,
		})

		f, err := fs.Open("test/with-segment/test.txt")

		Expect(err).NotTo(HaveOccurred())
		data, err := io.ReadAll(f)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).To(Equal("testing"))
	})

	It("should stat a file in the nested fs", func() {
		testFs := afero.NewMemMapFs()
		err := afero.WriteFile(testFs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		fs := mapped.NewFs(map[string]afero.Fs{
			"test": testFs,
		})

		fi, err := fs.Stat("test/test.txt")

		Expect(err).NotTo(HaveOccurred())
		Expect(fi.Name()).To(Equal("test.txt"))
	})

	It("should match rooted paths", func() {
		testFs := afero.NewMemMapFs()
		err := afero.WriteFile(testFs, "/test.txt", []byte(""), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		fs := mapped.NewFs(map[string]afero.Fs{
			"test": testFs,
		})

		_, err = fs.Open("/test/test.txt")

		Expect(err).NotTo(HaveOccurred())
	})

	It("should write to the mapped fs", func() {
		testFs := afero.NewMemMapFs()
		fs := mapped.NewFs(map[string]afero.Fs{
			"test": testFs,
		})

		err := afero.WriteFile(fs, "test/test.txt", []byte("testing"), os.ModePerm)

		Expect(err).NotTo(HaveOccurred())
		f, err := testFs.Open("test.txt")
		Expect(err).NotTo(HaveOccurred())
		data, err := io.ReadAll(f)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).To(Equal("testing"))
	})

	It("should write to a single fs", func() {
		testFs1 := afero.NewMemMapFs()
		testFs2 := afero.NewMemMapFs()
		fs := mapped.NewFs(map[string]afero.Fs{
			"test1": testFs1,
			"test2": testFs2,
		})

		err := afero.WriteFile(fs, "test1/test.txt", []byte("testing"), os.ModePerm)

		Expect(err).NotTo(HaveOccurred())
		_, err = testFs1.Open("test.txt")
		Expect(err).NotTo(HaveOccurred())

		_, err = testFs2.Open("test.txt")
		Expect(err).To(MatchError(os.ErrNotExist))
		_, err = testFs2.Open("test1/test.txt")
		Expect(err).To(MatchError(os.ErrNotExist))
		_, err = testFs2.Open("test2/test.txt")
		Expect(err).To(MatchError(os.ErrNotExist))
	})

	It("should error when a mapped fs is not found", func() {
		testFs := afero.NewMemMapFs()
		err := afero.WriteFile(testFs, "test.txt", []byte("testing"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		fs := mapped.NewFs(map[string]afero.Fs{
			"test": testFs,
		})

		_, err = fs.Open("other/test.txt")

		Expect(err).To(MatchError(os.ErrNotExist))
	})
})
