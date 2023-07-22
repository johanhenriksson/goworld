package physics

import (
	"runtime"
	"unsafe"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/quat"
)

type Compound struct {
	shapeBase
	object.Component

	shapes []Shape
}

func NewCompound() *Compound {
	cmp := object.NewComponent(&Compound{
		shapeBase: newShapeBase(CompoundShape),
	})

	cmp.handle = shape_new_compound(unsafe.Pointer(cmp))

	runtime.SetFinalizer(cmp, func(c *Compound) {
		c.destroy()
	})

	return cmp
}

func (c *Compound) OnEnable() {
	c.shapes = object.GetAllInChildren[Shape](c)
	for _, shape := range c.shapes {
		pos := c.Transform().Unproject(shape.Transform().WorldPosition())
		rot := quat.Ident()
		compound_add_child(c.handle, shape.shape(), pos, rot)
	}
	c.OnChange().Emit(c)
}

func (c *Compound) destroy() {
	c.shapes = nil
	if c.handle != nil {
		shape_delete(&c.handle)
	}
}
