package layout_test

import (
	. "github.com/johanhenriksson/goworld/gui/layout"

	"github.com/johanhenriksson/goworld/gui/dimension"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/math/vec2"

	"testing"
)

func TestColumnLayoutFixed(t *testing.T) {
	a := rect.Create("a", &rect.Props{})
	b := rect.Create("b", &rect.Props{})
	parent := rect.Create("parent", &rect.Props{
		Width:  dimension.Fixed(120),
		Height: dimension.Fixed(70),
		Layout: Column{
			Padding: 10,
			Gutter:  10,
		},
	})
	parent.SetChildren([]widget.T{a, b})

	parent.Arrange(vec2.New(1000, 300))

	// inner bounds should be 100x50
	// gutter is 10px -> usable space 100x40

	assertDimensions(t, parent, vec2.New(0, 0), vec2.New(120, 70))
	assertDimensions(t, a, vec2.New(10, 10), vec2.New(100, 20))
	assertDimensions(t, b, vec2.New(10, 40), vec2.New(100, 20))
}

func TestColumnLayoutPercent(t *testing.T) {
	a := rect.Create("a", &rect.Props{
		Height: dimension.Percent(50),
	})
	b := rect.Create("b", &rect.Props{
		Height: dimension.Percent(30),
	})
	parent := rect.Create("parent", &rect.Props{})
	parent.SetChildren([]widget.T{a, b})

	parent.Arrange(vec2.New(100, 100))

	// inner bounds should be 100x50
	// gutter is 10px -> usable space 100x40

	assertDimensions(t, parent, vec2.New(0, 0), vec2.New(100, 80))
	assertDimensions(t, a, vec2.New(0, 0), vec2.New(100, 50))
	assertDimensions(t, b, vec2.New(0, 50), vec2.New(100, 30))
}

func TestColumnLayoutAutoHeight(t *testing.T) {
	a := rect.Create("a", &rect.Props{
		Width:  dimension.Fixed(10),
		Height: dimension.Fixed(10),
	})
	b := rect.Create("b", &rect.Props{
		Width:  dimension.Fixed(20),
		Height: dimension.Fixed(10),
	})
	parent := rect.Create("parent", &rect.Props{
		Layout: Column{
			Padding: 10,
			Gutter:  10,
		},
	})
	parent.SetChildren([]widget.T{a, b})

	parent.Arrange(vec2.New(1000, 300))

	// inner bounds should be 100x50
	// gutter is 10px -> usable space 100x40

	assertDimensions(t, parent, vec2.New(0, 0), vec2.New(40, 50))
	assertDimensions(t, a, vec2.New(10, 10), vec2.New(10, 10))
	assertDimensions(t, b, vec2.New(10, 30), vec2.New(20, 10))
}

func TestColumnLayoutMixed(t *testing.T) {
	a := rect.Create("a", &rect.Props{
		Height: dimension.Fixed(10),
	})
	b := rect.Create("b", &rect.Props{})
	c := rect.Create("c", &rect.Props{})
	parent := rect.Create("parent", &rect.Props{
		Width:  dimension.Fixed(100),
		Height: dimension.Fixed(100),
		Layout: Column{
			Padding: 10,
			Gutter:  10,
		},
	})
	parent.SetChildren([]widget.T{a, b, c})

	parent.Arrange(vec2.New(1000, 300))

	// inner bounds should be 80x80
	// 2x gutter 10px -> usable space 80x60

	assertDimensions(t, parent, vec2.New(0, 0), vec2.New(100, 100))
	assertDimensions(t, a, vec2.New(10, 10), vec2.New(80, 10))

	// bottom elements share the remaining 50px height
	assertDimensions(t, b, vec2.New(10, 30), vec2.New(80, 25))
	assertDimensions(t, c, vec2.New(10, 65), vec2.New(80, 25))
}

func TestColumnLayoutFixedChildren(t *testing.T) {
	a := rect.Create("a", &rect.Props{
		Height: dimension.Fixed(30),
	})
	b := rect.Create("b", &rect.Props{
		Height: dimension.Fixed(30),
	})
	parent := rect.Create("parent", &rect.Props{
		Layout: Column{},
	})
	parent.SetChildren([]widget.T{a, b})

	parent.Arrange(vec2.New(100, 300))

	assertDimensions(t, parent, vec2.New(0, 0), vec2.New(100, 60))
	assertDimensions(t, a, vec2.New(0, 0), vec2.New(100, 30))
	assertDimensions(t, b, vec2.New(0, 30), vec2.New(100, 30))
}
func TestColumnLayoutWeird(t *testing.T) {
	a := rect.Create("a", &rect.Props{
		Height: dimension.Fixed(30),
	})
	b := rect.Create("b", &rect.Props{
		Height: dimension.Percent(50),
	})
	c := rect.Create("c", &rect.Props{
		Height: dimension.Auto(),
	})

	parent := rect.Create("parent", &rect.Props{
		Layout: Column{},
	})
	parent.SetChildren([]widget.T{a, b, c})

	parent.Arrange(vec2.New(100, 300))

	assertDimensions(t, parent, vec2.New(0, 0), vec2.New(100, 300))
	assertDimensions(t, a, vec2.New(0, 0), vec2.New(100, 30))
	assertDimensions(t, b, vec2.New(0, 30), vec2.New(100, 135))
	assertDimensions(t, c, vec2.New(0, 165), vec2.New(100, 135))
}
