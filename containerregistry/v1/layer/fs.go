package layer

import (
	"archive/tar"
	"bytes"
	"io"
	"os"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/spf13/afero"
	"github.com/spf13/afero/tarfs"
)

// The afero version of this upstream issue for fs.FS
// https://github.com/google/go-containerregistry/issues/921#issuecomment-769252935

func FromFs(fs afero.Fs) (v1.Layer, error) {
	buf := &bytes.Buffer{}
	w := tar.NewWriter(buf)

	err := afero.Walk(fs, "/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if err := w.WriteHeader(&tar.Header{
			Name: path,
			Size: info.Size(),
		}); err != nil {
			return err
		}

		f, err := fs.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := io.Copy(w, f); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}

	return tarball.LayerFromOpener(func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBuffer(buf.Bytes())), nil
	})
}

func ToFs(l v1.Layer) (afero.Fs, error) {
	rc, err := l.Uncompressed()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	return tarfs.New(tar.NewReader(rc)), nil
}
