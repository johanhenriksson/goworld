package style_test

import (
	. "github.com/johanhenriksson/goworld/gui/style"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/math/vec2"

	"github.com/kjk/flex"
)

var _ = Describe("column layout", func() {
	It("correctly sizes elements in a column layout", func() {

		a := rect.Create("a", rect.Props{
			Style: rect.Style{
				Height: Pct(50),
			},
		})
		b := rect.Create("b", rect.Props{
			Style: rect.Style{
				Height: Pct(50),
			},
		})
		parent := rect.Create("parent", rect.Props{
			Style: rect.Style{
				Layout: Column{},
			},
		})
		parent.SetChildren([]widget.T{a, b})
		root := parent.Flex()
		flex.CalculateLayout(root, 100, 100, flex.DirectionLTR)

		Expect(parent.Size()).To(Equal(vec2.New(100, 100)))
		Expect(a.Size()).To(Equal(vec2.New(100, 50)))
		Expect(a.Position()).To(Equal(vec2.New(0, 0)))
		Expect(b.Size()).To(Equal(vec2.New(100, 50)))
		Expect(b.Position()).To(Equal(vec2.New(0, 50)))
	})
})
