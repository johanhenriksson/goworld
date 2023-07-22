package physics

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec3"
)

//
// Box shape
//

type Box struct {
	shapeBase
	object.Component

	Extents *object.Property[vec3.T]
}

var _ Shape = &Box{}

func NewBox(size vec3.T) *Box {
	box := object.NewComponent(&Box{
		shapeBase: newShapeBase(BoxShape),
		Extents:   object.NewProperty(size),
	})

	// resize shape when extents are modified
	box.Extents.OnChange.Subscribe(box, box.resize)

	// trigger initial resize
	box.resize(size)

	runtime.SetFinalizer(box, func(b *Box) {
		b.destroy()
	})
	return box
}

func (b *Box) resize(size vec3.T) {
	b.destroy()
	b.handle = shape_new_box(unsafe.Pointer(b), size)
	b.OnChange().Emit(b)
}

func (b *Box) destroy() {
	if b.handle != nil {
		shape_delete(&b.handle)
	}
}

func (b *Box) Name() string {
	return "BoxCollider"
}

func (b *Box) String() string {
	return fmt.Sprintf("Box[Size=%s]", b.Extents.Get())
}
