package physics

import (
	"fmt"
	"log"
	"runtime"
	"unsafe"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/transform"
	"github.com/johanhenriksson/goworld/math/vec3"
)

//
// Box shape
//

type Box struct {
	shapeBase
	object.Component

	compound bool
	Extents  *object.Property[vec3.T]

	unsubTf func()
}

var _ Shape = &Box{}

func NewBox(size vec3.T) *Box {
	box := object.NewComponent(&Box{
		shapeBase: newShapeBase(BoxShape),
		Extents:   object.NewProperty(size),
	})

	// resize shape when extents are modified
	box.Extents.OnChange.Subscribe(box.resize)

	// trigger initial resize
	box.resize(size)

	runtime.SetFinalizer(box, func(b *Box) {
		b.destroy()
	})
	return box
}

func (b *Box) scale() vec3.T {
	if b.compound {
		return b.Transform().Scale()
	}
	return b.Transform().WorldScale()
}

func (b *Box) OnEnable() {
	// check if the box is part of a compound shape.
	// if it is, it should be scaled according to its local scale factor
	// otherwise, use world scale
	b.compound = hasParentShape(b.Parent()) || hasParentShape(b.Parent().Parent())
	log.Println("box", b.Parent().Name(), "is compound:", b.compound)

	// react to scale changes
	lastScale := b.scale()
	shape_scaling_set(b.handle, lastScale)
	b.OnChange().Emit(b)

	b.unsubTf = b.Transform().OnChange().Subscribe(func(t transform.T) {
		newScale := b.scale()
		if newScale != lastScale {
			log.Println("box scale update", b.Parent().Name(), ":", newScale)
			lastScale = newScale
			shape_scaling_set(b.handle, newScale)
			b.OnChange().Emit(b)
			// raising OnChange is technically not required since we dont recreate the shape
		}
	})
}

func (b *Box) OnDisable() {
	b.unsubTf()
}

func (b *Box) resize(size vec3.T) {
	b.destroy()
	b.handle = shape_new_box(unsafe.Pointer(b), size)
	shape_scaling_set(b.handle, b.scale())
	b.OnChange().Emit(b)
}

func (b *Box) destroy() {
	if b.handle != nil {
		shape_delete(&b.handle)
	}
}

func (b *Box) Name() string {
	return "BoxShape"
}

func (b *Box) String() string {
	return fmt.Sprintf("Box[Size=%s]", b.Extents.Get())
}
