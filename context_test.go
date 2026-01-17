package aferox_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/spf13/afero"
	"github.com/unmango/aferox"
)

// Mock setter for testing
type mockSetter struct {
	afero.Fs
	ctx context.Context
}

func (m *mockSetter) SetContext(ctx context.Context) {
	m.ctx = ctx
}

func (m *mockSetter) Name() string {
	return "MockSetterFs"
}

var _ = Describe("Context", func() {
	Describe("SetContext", func() {
		It("should set context on a context-supporting Fs", func() {
			ctx := context.Background()
			mock := &mockSetter{Fs: afero.NewMemMapFs()}

			err := aferox.SetContext(mock, ctx)

			Expect(err).NotTo(HaveOccurred())
			Expect(mock.ctx).To(BeIdenticalTo(ctx))
		})

		It("should return error when Fs does not support context", func() {
			ctx := context.Background()
			fs := afero.NewMemMapFs()

			err := aferox.SetContext(fs, ctx)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("context not supported"))
		})
	})

	Describe("FromContext", func() {
		It("should return default Fs when context has no Fs", func() {
			ctx := context.Background()

			fs := aferox.FromContext(ctx)

			Expect(fs).NotTo(BeNil())
			Expect(fs.Name()).To(Equal("OsFs"))
		})

		It("should return default Fs when context value is nil", func() {
			// Create a context with a different key
			ctx := context.WithValue(context.Background(), "different-key", nil)

			fs := aferox.FromContext(ctx)

			Expect(fs).NotTo(BeNil())
			Expect(fs.Name()).To(Equal("OsFs"))
		})
	})
})
