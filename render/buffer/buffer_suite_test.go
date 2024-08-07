package buffer_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
)

func TestBuffer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "render/buffer")
}

type mockBuffer struct {
	size int
}

var _ buffer.T = (*mockBuffer)(nil)

func (d *mockBuffer) Size() int                      { return d.size }
func (d *mockBuffer) Ptr() core1_0.Buffer            { return nil }
func (d *mockBuffer) Memory() device.Memory          { return nil }
func (d *mockBuffer) Read(offset int, data any) int  { return 0 }
func (d *mockBuffer) Write(offset int, data any) int { return 0 }

func (d *mockBuffer) Flush()   {}
func (d *mockBuffer) Destroy() {}
