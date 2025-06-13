package protofsv1alpha1_test

import (
	"fmt"
	"io"
	"net"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/spf13/afero"
	protofsv1alpha1 "github.com/unmango/aferox/protofs/grpc/v1alpha1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ = Describe("E2e", func() {
	var client, fs afero.Fs

	BeforeEach(func() {
		fs = afero.NewMemMapFs()
		server := grpc.NewServer()
		protofsv1alpha1.RegisterFsServer(server, fs)
		protofsv1alpha1.RegisterFileServer(server, fs)

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

	It("should create a writable file", func() {
		file, err := client.Create("test.txt")
		Expect(err).NotTo(HaveOccurred())

		_, err = io.WriteString(file, "testing text")
		Expect(err).NotTo(HaveOccurred())

		actual, err := fs.Open("test.txt")
		Expect(err).NotTo(HaveOccurred())
		data, err := io.ReadAll(actual)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).To(Equal("testing text"))
	})
})
