package color_test

import (
	"testing"

	"github.com/johanhenriksson/goworld/render/color"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestColor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "render/color")
}

var _ = Describe("colors", func() {
	It("converts from hex codes", func() {
		c := color.Hex("#123456")
		Expect(c.R).To(BeNumerically("~", float32(0x12)/255.0))
		Expect(c.G).To(BeNumerically("~", float32(0x34)/255.0))
		Expect(c.B).To(BeNumerically("~", float32(0x56)/255.0))

		a := color.Hex("#000000f0")
		Expect(a.A).To(BeNumerically("~", float32(0xF0)/255.0))
	})

	It("converts to hex codes", func() {
		c := color.RGB(
			float32(0x12)/255.0,
			float32(0x34)/255.0,
			float32(0x56)/255.0,
		)
		Expect(c.Hex()).To(Equal("#123456"))

		a := color.RGBA(
			0, 0, 0,
			float32(0xF0)/255.0,
		)
		Expect(a.Hex()).To(Equal("#000000f0"))
	})
})
