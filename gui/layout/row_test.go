package layout_test

import (
	"github.com/johanhenriksson/goworld/gui/dimension"
	. "github.com/johanhenriksson/goworld/gui/layout"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/math/vec2"

	"testing"
)

func TestRowLayoutFixed(t *testing.T) {
	a := rect.Create("a", &rect.Props{
		Width:  dimension.Auto(),
		Height: dimension.Auto(),
		Layout: Absolute{},
	})
	b := rect.Create("b", &rect.Props{
		Width:  dimension.Auto(),
		Height: dimension.Auto(),
		Layout: Absolute{},
	})
	parent := rect.Create("parent", &rect.Props{
		Width:  dimension.Fixed(100),
		Height: dimension.Fixed(100),
		Layout: Row{
			Padding: 10,
			Gutter:  10,
		},
	})
	parent.SetChildren([]widget.T{a, b})

	parent.Arrange(vec2.New(1000, 300))

	// inner bounds should be 80x80
	// gutter is 10px -> usable space 70x80

	assertDimensions(t, parent, vec2.New(0, 0), vec2.New(100, 100))
	assertDimensions(t, a, vec2.New(10, 10), vec2.New(35, 80))
	assertDimensions(t, b, vec2.New(55, 10), vec2.New(35, 80))
}

func TestRowLayoutAuto(t *testing.T) {
	a := rect.Create("a", &rect.Props{
		Width:  dimension.Fixed(10),
		Height: dimension.Fixed(10),
		Layout: Absolute{},
	})
	b := rect.Create("b", &rect.Props{
		Width:  dimension.Fixed(10),
		Height: dimension.Fixed(20),
		Layout: Absolute{},
	})
	parent := rect.Create("parent", &rect.Props{
		Width:  dimension.Auto(),
		Height: dimension.Auto(),
		Layout: Row{
			Padding: 10,
			Gutter:  10,
		},
	})
	parent.SetChildren([]widget.T{a, b})

	parent.Arrange(vec2.New(1000, 300))

	// inner bounds should be 80x80
	// gutter is 10px -> usable space 70x80

	assertDimensions(t, parent, vec2.New(0, 0), vec2.New(50, 40))
	assertDimensions(t, a, vec2.New(10, 10), vec2.New(10, 10))
	assertDimensions(t, b, vec2.New(30, 10), vec2.New(10, 20))
}

func TestRowLayoutMixed(t *testing.T) {
	a := rect.Create("a", &rect.Props{
		Width:  dimension.Fixed(10),
		Height: dimension.Auto(),
		Layout: Absolute{},
	})
	b := rect.Create("b", &rect.Props{
		Width:  dimension.Auto(),
		Height: dimension.Fixed(20),
		Layout: Absolute{},
	})
	parent := rect.Create("parent", &rect.Props{
		Width:  dimension.Fixed(100),
		Height: dimension.Auto(),
		Layout: Row{
			Padding: 10,
			Gutter:  10,
		},
	})
	parent.SetChildren([]widget.T{a, b})

	parent.Arrange(vec2.New(1000, 50))

	// inner bounds should be 80x80
	// gutter is 10px -> usable space 70x80

	assertDimensions(t, parent, vec2.New(0, 0), vec2.New(100, 50))
	assertDimensions(t, a, vec2.New(10, 10), vec2.New(10, 30))
	assertDimensions(t, b, vec2.New(30, 10), vec2.New(60, 20))
}
