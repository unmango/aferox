package grpc

import (
	"context"

	"buf.build/gen/go/unmango/protofs/grpc/go/dev/unmango/fs/v1alpha1/fsv1alpha1grpc"
	fsv1alpha1 "buf.build/gen/go/unmango/protofs/protocolbuffers/go/dev/unmango/fs/v1alpha1"
	"github.com/spf13/afero"
	"github.com/unmango/aferox/protofs/internal"
)

type Server struct {
	fsv1alpha1grpc.UnimplementedFsServiceServer
	Fs afero.Fs
}

func (s *Server) Chmod(_ context.Context, req *fsv1alpha1.ChmodRequest) (*fsv1alpha1.ChmodResponse, error) {
	if err := s.Fs.Chmod(req.Name, internal.OsFileMode(req.Mode)); err != nil {
		return &fsv1alpha1.ChmodResponse{Error: err.Error()}, nil
	} else {
		return &fsv1alpha1.ChmodResponse{}, nil
	}
}

func (s *Server) Chown(_ context.Context, req *fsv1alpha1.ChownRequest) (*fsv1alpha1.ChownResponse, error) {
	if err := s.Fs.Chown(req.Name, int(req.Uid), int(req.Gid)); err != nil {
		return &fsv1alpha1.ChownResponse{Error: err.Error()}, nil
	} else {
		return &fsv1alpha1.ChownResponse{}, nil
	}
}

func (s *Server) Chtimes(_ context.Context, req *fsv1alpha1.ChtimesRequest) (*fsv1alpha1.ChtimesResponse, error) {
	if err := s.Fs.Chtimes(req.Name, req.Atime.AsTime(), req.Mtime.AsTime()); err != nil {
		return &fsv1alpha1.ChtimesResponse{Error: err.Error()}, nil
	} else {
		return &fsv1alpha1.ChtimesResponse{}, nil
	}
}

func (s *Server) Create(_ context.Context, req *fsv1alpha1.CreateRequest) (*fsv1alpha1.CreateResponse, error) {
	if _, err := s.Fs.Create(req.Name); err != nil {
		return &fsv1alpha1.CreateResponse{Error: err.Error()}, nil
	} else {
		return &fsv1alpha1.CreateResponse{}, nil
	}
}
