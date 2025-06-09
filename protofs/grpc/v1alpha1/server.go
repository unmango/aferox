package protofsv1alpha1

import (
	"context"

	"buf.build/gen/go/unmango/protofs/grpc/go/dev/unmango/fs/v1alpha1/fsv1alpha1grpc"
	fsv1alpha1 "buf.build/gen/go/unmango/protofs/protocolbuffers/go/dev/unmango/fs/v1alpha1"
	"github.com/spf13/afero"
	"github.com/unmango/aferox/protofs/internal"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	fsv1alpha1grpc.UnimplementedFsServiceServer

	Fs afero.Fs
}

func (s *Server) Chmod(_ context.Context, req *fsv1alpha1.ChmodRequest) (*fsv1alpha1.ChmodResponse, error) {
	if err := s.Fs.Chmod(req.Name, internal.OsFileMode(req.Mode)); err != nil {
		return nil, err
	} else {
		return &fsv1alpha1.ChmodResponse{}, nil
	}
}

func (s *Server) Chown(_ context.Context, req *fsv1alpha1.ChownRequest) (*fsv1alpha1.ChownResponse, error) {
	if err := s.Fs.Chown(req.Name, int(req.Uid), int(req.Gid)); err != nil {
		return nil, err
	} else {
		return &fsv1alpha1.ChownResponse{}, nil
	}
}

func (s *Server) Chtimes(_ context.Context, req *fsv1alpha1.ChtimesRequest) (*fsv1alpha1.ChtimesResponse, error) {
	if err := s.Fs.Chtimes(req.Name, req.Atime.AsTime(), req.Mtime.AsTime()); err != nil {
		return nil, err
	} else {
		return &fsv1alpha1.ChtimesResponse{}, nil
	}
}

func (s *Server) Create(_ context.Context, req *fsv1alpha1.CreateRequest) (*fsv1alpha1.CreateResponse, error) {
	if file, err := s.Fs.Create(req.Name); err != nil {
		return nil, err
	} else {
		return &fsv1alpha1.CreateResponse{
			File: &fsv1alpha1.File{
				Name: file.Name(),
			},
		}, nil
	}
}

func (s *Server) Mkdir(_ context.Context, req *fsv1alpha1.MkdirRequest) (*fsv1alpha1.MkdirResponse, error) {
	if err := s.Fs.Mkdir(req.Name, internal.OsFileMode(req.Perm)); err != nil {
		return nil, err
	} else {
		return &fsv1alpha1.MkdirResponse{}, nil
	}
}

func (s *Server) Open(_ context.Context, req *fsv1alpha1.OpenRequest) (*fsv1alpha1.OpenResponse, error) {
	if file, err := s.Fs.Open(req.Name); err != nil {
		return nil, err
	} else {
		return &fsv1alpha1.OpenResponse{
			File: &fsv1alpha1.File{
				Name: file.Name(),
			},
		}, nil
	}
}

func (s *Server) OpenFile(_ context.Context, req *fsv1alpha1.OpenFileRequest) (*fsv1alpha1.OpenFileResponse, error) {
	if file, err := s.Fs.OpenFile(req.Name, int(req.Flag), internal.OsFileMode(req.Perm)); err != nil {
		return nil, err
	} else {
		return &fsv1alpha1.OpenFileResponse{
			File: &fsv1alpha1.File{
				Name: file.Name(),
			},
		}, nil
	}
}

func (s *Server) Remove(_ context.Context, req *fsv1alpha1.RemoveRequest) (*fsv1alpha1.RemoveResponse, error) {
	if err := s.Fs.Remove(req.Name); err != nil {
		return nil, err
	} else {
		return &fsv1alpha1.RemoveResponse{}, nil
	}
}

func (s *Server) RemoveAll(_ context.Context, req *fsv1alpha1.RemoveAllRequest) (*fsv1alpha1.RemoveAllResponse, error) {
	if err := s.Fs.RemoveAll(req.Path); err != nil {
		return nil, err
	} else {
		return &fsv1alpha1.RemoveAllResponse{}, nil
	}
}

func (s *Server) Rename(_ context.Context, req *fsv1alpha1.RenameRequest) (*fsv1alpha1.RenameResponse, error) {
	if err := s.Fs.Rename(req.Oldname, req.Newname); err != nil {
		return nil, err
	} else {
		return &fsv1alpha1.RenameResponse{}, nil
	}
}

func (s *Server) Stat(_ context.Context, req *fsv1alpha1.StatRequest) (*fsv1alpha1.StatResponse, error) {
	if info, err := s.Fs.Stat(req.Name); err != nil {
		return nil, err
	} else {
		return &fsv1alpha1.StatResponse{
			FileInfo: &fsv1alpha1.FileInfo{
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
