package protofsv1alpha1_test

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"buf.build/gen/go/unmango/protofs/grpc/go/dev/unmango/fs/v1alpha1/fsv1alpha1grpc"
	fsv1alpha1 "buf.build/gen/go/unmango/protofs/protocolbuffers/go/dev/unmango/fs/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/utils/ptr"

	"github.com/spf13/afero"
	protofsv1alpha1 "github.com/unmango/aferox/protofs/grpc/v1alpha1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ = Describe("Fs", func() {
	var (
		fs     afero.Fs
		client fsv1alpha1grpc.FsServiceClient
	)

	BeforeEach(func() {
		fs = afero.NewMemMapFs()
		server := grpc.NewServer()
		protofsv1alpha1.RegisterFsServer(server, fs)

		tmp := GinkgoT().TempDir()
		sock := filepath.Join(tmp, "fs.sock")

		lis, err := net.Listen("unix", sock)
		Expect(err).NotTo(HaveOccurred())

		By("Starting the FS server")
		go server.Serve(lis)
		DeferCleanup(server.GracefulStop)

		By("Creating the FS client")
		conn, err := grpc.NewClient(fmt.Sprint("unix://", sock),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).NotTo(HaveOccurred())
		client = fsv1alpha1grpc.NewFsServiceClient(conn)
	})

	It("should create a file", func(ctx context.Context) {
		file, err := client.Create(ctx, &fsv1alpha1.CreateRequest{
			Name: "test.txt",
		})
		Expect(err).NotTo(HaveOccurred())

		Expect(file.File.Name).To(Equal("test.txt"))
		Expect(file.File.Flag).To(Equal(ptr.To(int64(os.O_CREATE | os.O_RDWR))))
		Expect(file.File.Perm).To(BeNil())

		stat, err := fs.Stat("test.txt")
		Expect(err).NotTo(HaveOccurred())
		Expect(stat.Name()).To(Equal("test.txt"))
	})

	It("should open an existing file", func(ctx context.Context) {
		_, err := fs.Create("test.txt")
		Expect(err).NotTo(HaveOccurred())

		file, err := client.Open(ctx, &fsv1alpha1.OpenRequest{
			Name: "test.txt",
		})

		Expect(err).NotTo(HaveOccurred())
		Expect(file.File.Name).To(Equal("test.txt"))
		Expect(file.File.Flag).To(BeNil())
		Expect(file.File.Perm).To(BeNil())
	})

	It("should create a directory", func(ctx context.Context) {
		_, err := client.Mkdir(ctx, &fsv1alpha1.MkdirRequest{
			Name: "testdir",
		})
		Expect(err).NotTo(HaveOccurred())

		stat, err := fs.Stat("testdir")
		Expect(err).NotTo(HaveOccurred())
		Expect(stat.Name()).To(Equal("testdir"))
		Expect(stat.IsDir()).To(BeTrueBecause("file is a directory"))
	})
})
