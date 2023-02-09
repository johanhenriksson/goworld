package test

import (
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/kjk/flex"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("layout", func() {
	It("lays out flex rows properly", func() {
		small := rect.New("small", rect.Props{
			Style: rect.Style{
				Grow:   style.Grow(0),
				Shrink: style.Shrink(0),
				Basis:  style.Px(15),
			},
		})
		big := rect.New("big", rect.Props{
			Style: rect.Style{
				Basis:  style.Px(100),
				Shrink: style.Shrink(1),
				Grow:   style.Grow(0),
			},
		})
		row := rect.New("row", rect.Props{
			Style: rect.Style{
				Layout: style.Row{},
				Grow:   style.Grow(0),
			},
			Children: []node.T{
				small, big,
			},
		})
		tree := row.Hydrate("root")
		flex.CalculateLayout(tree.Flex(), 100, 10, flex.DirectionLTR)

		Expect(tree.Position()).To(Equal(vec2.New(0, 0)))
		Expect(tree.Size()).To(Equal(vec2.New(100, 10)))

		Expect(tree.Children()[0].Position()).To(Equal(vec2.New(0, 0)))
		Expect(tree.Children()[0].Size()).To(Equal(vec2.New(15, 10)))

		Expect(tree.Children()[1].Position()).To(Equal(vec2.New(15, 0)))
		Expect(tree.Children()[1].Size()).To(Equal(vec2.New(85, 10)))
	})
})
