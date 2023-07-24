package physics

import (
	"fmt"
	"unsafe"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Box struct {
	kind ShapeType
	*Collider

	Extents *object.Property[vec3.T]
}

var _ = checkShape(NewBox(vec3.Zero))

func NewBox(size vec3.T) *Box {
	box := object.NewComponent(&Box{
		kind:    BoxShape,
		Extents: object.NewProperty(size),
	})
	box.Collider = newCollider(box)

	// resize shape when extents are modified
	box.Extents.OnChange.Subscribe(func(t vec3.T) {
		box.refresh()
	})

	return box
}

func (b *Box) Name() string {
	return "BoxShape"
}

func (b *Box) String() string {
	return fmt.Sprintf("Box[Size=%s]", b.Extents.Get())
}

func (b *Box) colliderCreate() shapeHandle {
	return shape_new_box(unsafe.Pointer(b), b.Extents.Get())
}

func (b *Box) colliderIsCompound() bool {
	return defaultCompoundCheck(b)
}

func (b *Box) colliderDestroy() {}
