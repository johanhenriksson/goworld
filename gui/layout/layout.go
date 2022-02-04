package layout

import (
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/vec2"
)

type T interface {
	Flow(Layoutable)
}

type Layoutable interface {
	Size() vec2.T
	Children() []widget.T
}

// we want to layout something.
// that imples that the something has child elements.
// the child elements may have their own children
// if not, they will be of a known (?) size
//   - no! it depends. eg. if a label has its width constrained,
//     we want to flow the text over multiple lines which increases
//     the height
//
// so we introduce 3 types of layouts
//   - absolute: no layouting. elements remain in their original position
//   - column: children inherit parent width, but may have a dynamic height
//   - row: children inherit parent height, but may have a dynamic width
//
// in other words, layouts may constrain elements in a single axis, but then
// they must be free to adjust in the other.
//
// elements should have a function that returns their desired size given
// a specific constraint. this can be computed bottom-up recursively
// actually, the elements could probably be resized on the way up the recursion
// the layouting would then end with resizing the original element
