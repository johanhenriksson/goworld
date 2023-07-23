package physics

import (
	"log"
	"runtime"
	"unsafe"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/transform"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Compound struct {
	shapeBase
	object.Component

	compound bool
	shapes   []*childShape
	unsubTf  func()
}

type childShape struct {
	Shape
	index    int
	localPos vec3.T
	localRot quat.T
	unsub    func()
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

func (c *Compound) Name() string {
	return "CompoundShape"
}

func (c *Compound) Update(scene object.Component, dt float32) {
	// log.Println("compound world scale", c.Transform().WorldScale().Scaled(1.2))
	// shape_scaling_set(c.handle, c.Transform().WorldScale())
}

func (c *Compound) scale() vec3.T {
	if c.compound {
		return c.Transform().Scale()
	}
	return c.Transform().WorldScale()
}

func (c *Compound) OnEnable() {
	c.compound = object.Get[*RigidBody](c) == nil && hasParentShape(c.Parent().Parent())
	c.refresh()

	// react to scale changes
	lastScale := c.scale()
	shape_scaling_set(c.handle, lastScale)
	c.OnChange().Emit(c)

	c.unsubTf = c.Transform().OnChange().Subscribe(func(t transform.T) {
		newScale := c.scale()
		if newScale != lastScale {
			lastScale = newScale
			shape_scaling_set(c.handle, newScale)
			c.OnChange().Emit(c)
			// raising OnChange is technically not required since we dont recreate the shape
		}
	})
}

func (c *Compound) refresh() {
	log.Println("refresh compound shape", c.Parent().Name())
	c.destroy()
	c.handle = shape_new_compound(unsafe.Pointer(c))

	// find all shapes that should be combined into the compound mesh
	// todo: also subscribe to attach/detach events on all relevant child objects
	shapes := Shapes(c)
	c.shapes = c.shapes[:0]
	for i, shape := range shapes {
		if shape.shape() == nil {
			panic("child shape is nil")
		}

		// find rotation relative to compound shape
		child := &childShape{
			Shape:    shape,
			index:    i,
			localPos: shape.Transform().Position().Mul(c.scale()),
			localRot: shape.Transform().Rotation(),
		}
		c.shapes = append(c.shapes, child)

		// shape_scaling_set(shape.shape(), localScale)
		compound_add_child(c.handle, shape.shape(), child.localPos, child.localRot)
		log.Println("add compound child", shape.Parent().Name(), child.localPos)

		// child shape changes should trigger a complete recreation of the compound shape
		unsubShape := shape.OnChange().Subscribe(func(s Shape) {
			c.refresh()
		})

		// adjust scale & local position on transform changes
		unsubTf := shape.Transform().OnChange().Subscribe(func(t transform.T) {
			newPos := t.Position().Mul(c.scale())
			newRot := t.Rotation()

			if !newPos.ApproxEqual(child.localPos) || !newRot.ApproxEqual(child.localRot) {
				compound_update_child(c.handle, child.index, newPos, newRot)
				child.localPos = newPos
				child.localRot = newRot
				c.OnChange().Emit(c)
			}
		})

		child.unsub = func() {
			unsubShape()
			unsubTf()
		}
	}
	c.OnChange().Emit(c)
}

func (c *Compound) destroy() {
	for _, shape := range c.shapes {
		shape.unsub()
	}
	c.shapes = nil
	if c.handle != nil {
		shape_delete(&c.handle)
	}
}
