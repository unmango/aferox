package protofsv1alpha1_test

import (
	"fmt"
	"net"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/spf13/afero"
	protofsv1alpha1 "github.com/unmango/aferox/protofs/grpc/v1alpha1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ = Describe("Fs", func() {
	var (
		client, fs afero.Fs
		server     *grpc.Server
	)

	BeforeEach(func() {
		fs = afero.NewMemMapFs()
		server = grpc.NewServer()
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
		client = protofsv1alpha1.NewFs(conn)
	})

	It("should create a file", func() {
		_, err := client.Create("test.txt")
		Expect(err).NotTo(HaveOccurred())
	})
})
