package physics

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/johanhenriksson/goworld/core/object"
)

type Capsule struct {
	shapeBase
	object.Component

	Radius *object.Property[float32]
	Height *object.Property[float32]
}

var _ = &Capsule{}

func NewCapsule(height, radius float32) *Capsule {
	capsule := object.NewComponent(&Capsule{
		shapeBase: newShapeBase(CapsuleShape),
		Radius:    object.NewProperty(radius),
		Height:    object.NewProperty(height),
	})

	capsule.Radius.OnChange.Subscribe(capsule, func(radius float32) {
		capsule.resize(radius, capsule.Height.Get())
	})
	capsule.Height.OnChange.Subscribe(capsule, func(height float32) {
		capsule.resize(capsule.Radius.Get(), height)
	})

	// trigger initial resize
	capsule.resize(radius, height)

	runtime.SetFinalizer(capsule, func(c *Capsule) {
		c.destroy()
	})
	return capsule
}

func (c *Capsule) resize(radius, height float32) {
	c.destroy()
	c.handle = shape_new_capsule(unsafe.Pointer(c), radius, height)
	c.OnChange().Emit(c)
}

func (c *Capsule) destroy() {
	if c.handle != nil {
		shape_delete(&c.handle)
	}
}

func (c *Capsule) String() string {
	return fmt.Sprintf("Capsule[Height=%.2f,Radius=%.2f]", c.Height.Get(), c.Radius.Get())
}
