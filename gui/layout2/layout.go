package layout2

import "github.com/johanhenriksson/goworld/math/vec2"

func HyperLayout() {
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
}

type Widget interface {
	Resize(vec2.T)
	Measure(vec2.T) vec2.T
	Children() []Widget
}

type Image struct {
	srcWidth  float32
	srcHeight float32
}

func (i Image) Measure(available vec2.T) vec2.T {
	// the image wants to maintain its aspect ratio
	aspect := i.srcWidth / i.srcHeight
	return vec2.New(available.X, available.X*aspect)
}

func (i Image) Resize(vec2.T)      {}
func (i Image) Children() []Widget { return nil }

type Frame struct {
}
