package physics

import (
	"fmt"
	"unsafe"

	"github.com/johanhenriksson/goworld/core/object"
)

type Capsule struct {
	kind ShapeType
	*Collider

	Radius *object.Property[float32]
	Height *object.Property[float32]
}

var _ = checkShape(NewCapsule(1, 1))

func NewCapsule(height, radius float32) *Capsule {
	capsule := object.NewComponent(&Capsule{
		kind:   CapsuleShape,
		Radius: object.NewProperty(radius),
		Height: object.NewProperty(height),
	})
	capsule.Collider = newCollider(capsule, true)

	capsule.Radius.OnChange.Subscribe(func(radius float32) {
		capsule.refresh()
	})
	capsule.Height.OnChange.Subscribe(func(height float32) {
		capsule.refresh()
	})

	return capsule
}

func (c *Capsule) Name() string {
	return "CapsuleShape"
}

func (c *Capsule) String() string {
	return fmt.Sprintf("Capsule[Height=%.2f,Radius=%.2f]", c.Height.Get(), c.Radius.Get())
}

func (c *Capsule) colliderCreate() shapeHandle {
	return shape_new_capsule(unsafe.Pointer(c), c.Radius.Get(), c.Height.Get())
}

func (c *Capsule) colliderRefresh() {}
func (c *Capsule) colliderDestroy() {}
