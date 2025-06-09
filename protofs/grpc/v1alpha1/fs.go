package protofsv1alpha1

import (
	"context"
	"os"
	"time"

	"buf.build/gen/go/unmango/protofs/grpc/go/dev/unmango/file/v1alpha1/filev1alpha1grpc"
	"buf.build/gen/go/unmango/protofs/grpc/go/dev/unmango/fs/v1alpha1/fsv1alpha1grpc"
	filev1alpha1 "buf.build/gen/go/unmango/protofs/protocolbuffers/go/dev/unmango/file/v1alpha1"
	fsv1alpha1 "buf.build/gen/go/unmango/protofs/protocolbuffers/go/dev/unmango/fs/v1alpha1"
	"github.com/spf13/afero"
	"github.com/unmango/aferox/protofs/internal"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Fs struct {
	conn   grpc.ClientConnInterface
	client fsv1alpha1grpc.FsServiceClient
}

func NewFs(conn grpc.ClientConnInterface) afero.Fs {
	return &Fs{
		conn:   conn,
		client: fsv1alpha1grpc.NewFsServiceClient(conn),
	}
}

// Chmod implements afero.Fs.
func (f *Fs) Chmod(name string, mode os.FileMode) error {
	_, err := f.client.Chmod(context.TODO(), &fsv1alpha1.ChmodRequest{
		Name: name,
		Mode: internal.ProtoFileMode(mode),
	})

	return err
}

// Chown implements afero.Fs.
func (f *Fs) Chown(name string, uid int, gid int) error {
	_, err := f.client.Chown(context.TODO(), &fsv1alpha1.ChownRequest{
		Name: name,
		Uid:  int32(uid),
		Gid:  int32(gid),
	})

	return err
}

// Chtimes implements afero.Fs.
func (f *Fs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	_, err := f.client.Chtimes(context.TODO(), &fsv1alpha1.ChtimesRequest{
		Name:  name,
		Atime: timestamppb.New(atime),
		Mtime: timestamppb.New(mtime),
	})

	return err
}

// Create implements afero.Fs.
func (f *Fs) Create(name string) (afero.File, error) {
	res, err := f.client.Create(context.TODO(), &fsv1alpha1.CreateRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}

	return File{
		client: filev1alpha1grpc.NewFileServiceClient(f.conn),
		file:   res.File,
	}, nil
}

// Mkdir implements afero.Fs.
func (f *Fs) Mkdir(name string, perm os.FileMode) error {
	_, err := f.client.Mkdir(context.TODO(), &fsv1alpha1.MkdirRequest{
		Name: name,
		Perm: internal.ProtoFileMode(perm),
	})

	return err
}

// MkdirAll implements afero.Fs.
func (f *Fs) MkdirAll(path string, perm os.FileMode) error {
	panic("unimplemented")
}

// Name implements afero.Fs.
func (f *Fs) Name() string {
	return "protofs"
}

// Open implements afero.Fs.
func (f *Fs) Open(name string) (afero.File, error) {
	res, err := f.client.Open(context.TODO(), &fsv1alpha1.OpenRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}

	return File{
		client: filev1alpha1grpc.NewFileServiceClient(f.conn),
		file:   res.File,
	}, nil
}

// OpenFile implements afero.Fs.
func (f *Fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	res, err := f.client.OpenFile(context.TODO(), &fsv1alpha1.OpenFileRequest{
		Name: name,
		Flag: int32(flag),
		Perm: internal.ProtoFileMode(perm),
	})
	if err != nil {
		return nil, err
	}

	return File{
		client: filev1alpha1grpc.NewFileServiceClient(f.conn),
		file:   res.File,
	}, nil
}

// Remove implements afero.Fs.
func (f *Fs) Remove(name string) error {
	_, err := f.client.Remove(context.TODO(), &fsv1alpha1.RemoveRequest{
		Name: name,
	})

	return err
}

// RemoveAll implements afero.Fs.
func (f *Fs) RemoveAll(path string) error {
	_, err := f.client.RemoveAll(context.TODO(), &fsv1alpha1.RemoveAllRequest{
		Path: path,
	})

	return err
}

// Rename implements afero.Fs.
func (f *Fs) Rename(oldname string, newname string) error {
	_, err := f.client.Rename(context.TODO(), &fsv1alpha1.RenameRequest{
		Oldname: oldname,
		Newname: newname,
	})

	return err
}

// Stat implements afero.Fs.
func (f *Fs) Stat(name string) (os.FileInfo, error) {
	res, err := f.client.Stat(context.TODO(), &fsv1alpha1.StatRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}

	return FileInfo{res.FileInfo}, nil
}

type FsServer struct {
	fsv1alpha1grpc.UnimplementedFsServiceServer

	Fs afero.Fs
}

func (s *FsServer) Chmod(_ context.Context, req *fsv1alpha1.ChmodRequest) (*fsv1alpha1.ChmodResponse, error) {
	if err := s.Fs.Chmod(req.Name, internal.OsFileMode(req.Mode)); err != nil {
		return nil, err
	} else {
		return &fsv1alpha1.ChmodResponse{}, nil
	}
}

func (s *FsServer) Chown(_ context.Context, req *fsv1alpha1.ChownRequest) (*fsv1alpha1.ChownResponse, error) {
	if err := s.Fs.Chown(req.Name, int(req.Uid), int(req.Gid)); err != nil {
		return nil, err
	} else {
		return &fsv1alpha1.ChownResponse{}, nil
	}
}

func (s *FsServer) Chtimes(_ context.Context, req *fsv1alpha1.ChtimesRequest) (*fsv1alpha1.ChtimesResponse, error) {
	if err := s.Fs.Chtimes(req.Name, req.Atime.AsTime(), req.Mtime.AsTime()); err != nil {
		return nil, err
	} else {
		return &fsv1alpha1.ChtimesResponse{}, nil
	}
}

func (s *FsServer) Create(_ context.Context, req *fsv1alpha1.CreateRequest) (*fsv1alpha1.CreateResponse, error) {
	if file, err := s.Fs.Create(req.Name); err != nil {
		return nil, err
	} else {
		return &fsv1alpha1.CreateResponse{
			File: &filev1alpha1.File{
				Name: file.Name(),
			},
		}, nil
	}
}

func (s *FsServer) Mkdir(_ context.Context, req *fsv1alpha1.MkdirRequest) (*fsv1alpha1.MkdirResponse, error) {
	if err := s.Fs.Mkdir(req.Name, internal.OsFileMode(req.Perm)); err != nil {
		return nil, err
	} else {
		return &fsv1alpha1.MkdirResponse{}, nil
	}
}

func (s *FsServer) Open(_ context.Context, req *fsv1alpha1.OpenRequest) (*fsv1alpha1.OpenResponse, error) {
	if file, err := s.Fs.Open(req.Name); err != nil {
		return nil, err
	} else {
		return &fsv1alpha1.OpenResponse{
			File: &filev1alpha1.File{
				Name: file.Name(),
			},
		}, nil
	}
}

func (s *FsServer) OpenFile(_ context.Context, req *fsv1alpha1.OpenFileRequest) (*fsv1alpha1.OpenFileResponse, error) {
	if file, err := s.Fs.OpenFile(req.Name, int(req.Flag), internal.OsFileMode(req.Perm)); err != nil {
		return nil, err
	} else {
		return &fsv1alpha1.OpenFileResponse{
			File: &filev1alpha1.File{
				Name: file.Name(),
			},
		}, nil
	}
}

func (s *FsServer) Remove(_ context.Context, req *fsv1alpha1.RemoveRequest) (*fsv1alpha1.RemoveResponse, error) {
	if err := s.Fs.Remove(req.Name); err != nil {
		return nil, err
	} else {
		return &fsv1alpha1.RemoveResponse{}, nil
	}
}

func (s *FsServer) RemoveAll(_ context.Context, req *fsv1alpha1.RemoveAllRequest) (*fsv1alpha1.RemoveAllResponse, error) {
	if err := s.Fs.RemoveAll(req.Path); err != nil {
		return nil, err
	} else {
		return &fsv1alpha1.RemoveAllResponse{}, nil
	}
}

func (s *FsServer) Rename(_ context.Context, req *fsv1alpha1.RenameRequest) (*fsv1alpha1.RenameResponse, error) {
	if err := s.Fs.Rename(req.Oldname, req.Newname); err != nil {
		return nil, err
	} else {
		return &fsv1alpha1.RenameResponse{}, nil
	}
}

func (s *FsServer) Stat(_ context.Context, req *fsv1alpha1.StatRequest) (*fsv1alpha1.StatResponse, error) {
	if info, err := s.Fs.Stat(req.Name); err != nil {
		return nil, err
	} else {
		return &fsv1alpha1.StatResponse{
			FileInfo: &filev1alpha1.FileInfo{
				Name:    info.Name(),
				Size:    info.Size(),
				Mode:    internal.ProtoFileMode(info.Mode()),
				ModTime: timestamppb.New(info.ModTime()),
				IsDir:   info.IsDir(),
				// Sys:  info.Sys(), // TODO
			},
		}, nil
	}
}

func RegisterFsServer(s grpc.ServiceRegistrar, fs afero.Fs) {
	fsv1alpha1grpc.RegisterFsServiceServer(s, &FsServer{Fs: fs})
}
