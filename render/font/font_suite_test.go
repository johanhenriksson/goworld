package font_test

import (
	"testing"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/render/font"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/image/math/fixed"
)

func TestFont(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "render/font")
}

var _ = Describe("font utils", func() {
	It("converts fixed to float32", func() {
		v := fixed.I(2)
		Expect(font.FixToFloat(v)).To(BeNumerically("~", float32(2.0)))

		v2 := fixed.Int26_6(1<<6 + 1<<4)
		Expect(font.FixToFloat(v2)).To(BeNumerically("~", float32(1.25)))
	})

	It("extracts glyphs", func() {
		f := font.Load(assets.FS, "fonts/SourceSansPro-Regular.ttf", 12, 1)
		Expect(f).ToNot(BeNil())
		a, err := f.Glyph('g')
		Expect(err).ToNot(HaveOccurred())
		Expect(a.Advance).To(BeNumerically(">", 0))
	})
})
