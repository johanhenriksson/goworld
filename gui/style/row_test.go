package style_test

import (
	"testing"

	. "github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/gui/widget/rect"

	"github.com/kjk/flex"
)

func TestRowFixedChildren(t *testing.T) {
	a := rect.Create("a", &rect.Props{
		Style: Sheet{
			Width:  Px(20),
			Height: Px(20),
		},
	})
	b := rect.Create("b", &rect.Props{
		Style: Sheet{
			Width:  Px(20),
			Height: Px(20),
		},
	})
	parent := rect.Create("parent", &rect.Props{
		Style: Sheet{
			Layout: Row{},
		},
	})
	parent.SetChildren([]widget.T{a, b})
	root := parent.Flex()
	flex.CalculateLayout(root, flex.Undefined, flex.Undefined, flex.DirectionLTR)

	assertSize(t, parent, 40, 20)
	assertSize(t, a, 20, 20)
	assertPosition(t, a, 0, 0)
	assertSize(t, b, 20, 20)
	assertPosition(t, b, 20, 0)
}
