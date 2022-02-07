package style_test

import (
	"testing"

	. "github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/gui/widget/rect"

	"github.com/kjk/flex"
)

func TestColumn(t *testing.T) {
	a := rect.Create("a", &rect.Props{
		Style: Sheet{
			Height: Pct(50),
		},
	})
	b := rect.Create("b", &rect.Props{
		Style: Sheet{
			Height: Pct(50),
		},
	})
	parent := rect.Create("parent", &rect.Props{
		Style: Sheet{
			Layout: Column{},
		},
	})
	parent.SetChildren([]widget.T{a, b})
	root := parent.Flex()
	flex.CalculateLayout(root, 100, 100, flex.DirectionLTR)

	assertSize(t, parent, 100, 100)
	assertSize(t, a, 100, 50)
	assertPosition(t, a, 0, 0)
	assertSize(t, b, 100, 50)
	assertPosition(t, b, 0, 50)
}
