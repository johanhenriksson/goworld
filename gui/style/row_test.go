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

var _ = Describe("row layout", func() {
	It("correctly sizes elements in a row layout", func() {
		a := rect.Create("a", rect.Props{
			Style: rect.Style{
				Width:  Px(20),
				Height: Px(20),
			},
		})
		b := rect.Create("b", rect.Props{
			Style: rect.Style{
				Width:  Px(20),
				Height: Px(20),
			},
		})
		parent := rect.Create("parent", rect.Props{
			Style: rect.Style{
				Layout: Row{},
			},
		})
		parent.SetChildren([]widget.T{a, b})
		root := parent.Flex()
		flex.CalculateLayout(root, flex.Undefined, flex.Undefined, flex.DirectionLTR)

		Expect(parent.Size()).To(Equal(vec2.New(40, 20)))
		Expect(a.Size()).To(Equal(vec2.New(20, 20)))
		Expect(b.Size()).To(Equal(vec2.New(20, 20)))
		Expect(a.Position()).To(Equal(vec2.New(0, 0)))
		Expect(b.Position()).To(Equal(vec2.New(20, 0)))
	})
})
