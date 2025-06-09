package protofsv1alpha1

import (
	"context"
	"os"
	"time"

	"buf.build/gen/go/unmango/protofs/grpc/go/dev/unmango/fs/v1alpha1/fsv1alpha1grpc"
	fsv1alpha1 "buf.build/gen/go/unmango/protofs/protocolbuffers/go/dev/unmango/fs/v1alpha1"
	"github.com/spf13/afero"
	"github.com/unmango/aferox/protofs/internal"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Fs struct {
	client fsv1alpha1grpc.FsServiceClient
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
	panic("unimplemented")
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
	panic("unimplemented")
}

// OpenFile implements afero.Fs.
func (f *Fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	panic("unimplemented")
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
	panic("unimplemented")
}

func New(conn grpc.ClientConnInterface) afero.Fs {
	client := fsv1alpha1grpc.NewFsServiceClient(conn)
	return &Fs{client: client}
}
