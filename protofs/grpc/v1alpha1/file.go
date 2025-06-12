package protofsv1alpha1

import (
	"context"
	"io"
	"os"

	"buf.build/gen/go/unmango/protofs/grpc/go/dev/unmango/file/v1alpha1/filev1alpha1grpc"
	filev1alpha1 "buf.build/gen/go/unmango/protofs/protocolbuffers/go/dev/unmango/file/v1alpha1"
	"github.com/spf13/afero"
	"github.com/unmango/aferox/protofs/internal"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type File struct {
	client filev1alpha1grpc.FileServiceClient
	file   *filev1alpha1.File
}

// Close implements afero.File.
func (f File) Close() error {
	return nil
}

// Name implements afero.File.
func (f File) Name() string {
	return f.file.Name
}

// Read implements afero.File.
func (f File) Read(p []byte) (n int, err error) {
	res, err := f.client.Read(context.TODO(), &filev1alpha1.ReadRequest{
		File: f.file,
	})
	if err != nil {
		return 0, err
	}

	return copy(res.Data, p), nil
}

// ReadAt implements afero.File.
func (f File) ReadAt(p []byte, off int64) (n int, err error) {
	res, err := f.client.ReadAt(context.Background(), &filev1alpha1.ReadAtRequest{
		File:   f.file,
		Offset: off,
	})
	if err != nil {
		return 0, err
	}

	return copy(res.Data, p), nil
}

// Readdir implements afero.File.
func (f File) Readdir(count int) (info []os.FileInfo, err error) {
	res, err := f.client.Readdir(context.TODO(), &filev1alpha1.ReaddirRequest{})
	if err != nil {
		return nil, err
	}

	for _, fi := range res.FileInfos {
		info = append(info, FileInfo{fi})
	}

	return info, nil
}

// Readdirnames implements afero.File.
func (f File) Readdirnames(n int) ([]string, error) {
	res, err := f.client.ReaddirNames(context.TODO(), &filev1alpha1.ReaddirNamesRequest{
		File:  f.file,
		Count: int32(n),
	})
	if err != nil {
		return nil, err
	}

	return res.Names, nil
}

// Seek implements afero.File.
func (f File) Seek(offset int64, whence int) (int64, error) {
	panic("unimplemented")
}

// Stat implements afero.File.
func (f File) Stat() (os.FileInfo, error) {
	res, err := f.client.Stat(context.TODO(), &filev1alpha1.StatRequest{
		File: f.file,
	})
	if err != nil {
		return nil, err
	}

	return FileInfo{res.FileInfo}, nil
}

// Sync implements afero.File.
func (f File) Sync() error {
	panic("unimplemented")
}

// Truncate implements afero.File.
func (f File) Truncate(size int64) error {
	_, err := f.client.Truncate(context.TODO(), &filev1alpha1.TruncateRequest{
		File: f.file,
		Size: size,
	})

	return err
}

// Write implements afero.File.
func (f File) Write(p []byte) (n int, err error) {
	_, err = f.client.Write(context.TODO(), &filev1alpha1.WriteRequest{
		File: f.file,
		Data: p,
	})
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

// WriteAt implements afero.File.
func (f File) WriteAt(p []byte, off int64) (n int, err error) {
	_, err = f.client.WriteAt(context.TODO(), &filev1alpha1.WriteAtRequest{
		File:   f.file,
		Data:   p,
		Offset: off,
	})
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

// WriteString implements afero.File.
func (f File) WriteString(s string) (ret int, err error) {
	return f.Write([]byte(s))
}

type FileServer struct {
	filev1alpha1grpc.UnimplementedFileServiceServer

	Fs afero.Fs
}

func (s *FileServer) Read(_ context.Context, req *filev1alpha1.ReadRequest) (*filev1alpha1.ReadResponse, error) {
	file, err := s.Fs.Open(req.File.Name)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return &filev1alpha1.ReadResponse{
		Data: data,
	}, nil
}

func (s *FileServer) ReadAt(_ context.Context, req *filev1alpha1.ReadAtRequest) (*filev1alpha1.ReadAtResponse, error) {
	file, err := s.Fs.Open(req.File.Name)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 0, info.Size())
	_, err = file.ReadAt(buf, req.Offset)
	if err != nil {
		return nil, err
	}

	return &filev1alpha1.ReadAtResponse{
		Data: buf,
	}, nil
}

func (s *FileServer) Readdir(_ context.Context, req *filev1alpha1.ReaddirRequest) (*filev1alpha1.ReaddirResponse, error) {
	file, err := s.Fs.Open(req.File.Name)
	if err != nil {
		return nil, err
	}

	info, err := file.Readdir(int(req.Count))
	if err != nil {
		return nil, err
	}

	res := &filev1alpha1.ReaddirResponse{}
	for _, fi := range info {
		res.FileInfos = append(res.FileInfos, &filev1alpha1.FileInfo{
			Name:    fi.Name(),
			Size:    fi.Size(),
			Mode:    internal.ProtoFileMode(fi.Mode()),
			ModTime: timestamppb.New(fi.ModTime()),
			IsDir:   fi.IsDir(),
		})
	}

	return res, nil
}

func (s *FileServer) ReaddirNames(_ context.Context, req *filev1alpha1.ReaddirNamesRequest) (*filev1alpha1.ReaddirNamesResponse, error) {
	file, err := s.Fs.Open(req.File.Name)
	if err != nil {
		return nil, err
	}

	names, err := file.Readdirnames(int(req.Count))
	if err != nil {
		return nil, err
	}

	return &filev1alpha1.ReaddirNamesResponse{
		Names: names,
	}, nil
}

func (s *FileServer) Stat(_ context.Context, req *filev1alpha1.StatRequest) (*filev1alpha1.StatResponse, error) {
	file, err := s.Fs.Open(req.File.Name)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	return &filev1alpha1.StatResponse{
		FileInfo: &filev1alpha1.FileInfo{
			Name:    info.Name(),
			Size:    info.Size(),
			Mode:    internal.ProtoFileMode(info.Mode()),
			ModTime: timestamppb.New(info.ModTime()),
			IsDir:   info.IsDir(),
		},
	}, nil
}

func (s *FileServer) Truncate(_ context.Context, req *filev1alpha1.TruncateRequest) (*filev1alpha1.TruncateResponse, error) {
	file, err := s.Fs.Open(req.File.Name)
	if err != nil {
		return nil, err
	}

	if err := file.Truncate(req.Size); err != nil {
		return nil, err
	} else {
		return &filev1alpha1.TruncateResponse{}, nil
	}
}

func (s *FileServer) Write(_ context.Context, req *filev1alpha1.WriteRequest) (*filev1alpha1.WriteResponse, error) {
	file, err := s.Fs.Open(req.File.Name)
	if err != nil {
		return nil, err
	}

	if _, err = file.Write(req.Data); err != nil {
		return nil, err
	} else {
		return &filev1alpha1.WriteResponse{}, nil
	}
}

func (s *FileServer) WriteAt(_ context.Context, req *filev1alpha1.WriteAtRequest) (*filev1alpha1.WriteAtResponse, error) {
	file, err := s.Fs.Open(req.File.Name)
	if err != nil {
		return nil, err
	}

	if _, err = file.WriteAt(req.Data, req.Offset); err != nil {
		return nil, err
	} else {
		return &filev1alpha1.WriteAtResponse{}, nil
	}
}

func RegisterFileServer(s grpc.ServiceRegistrar, fs afero.Fs) {
	filev1alpha1grpc.RegisterFileServiceServer(s, &FileServer{Fs: fs})
}
