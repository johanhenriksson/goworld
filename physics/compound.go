package physics

import (
	"unsafe"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/transform"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Compound struct {
	kind ShapeType
	*Collider

	shapes []*childShape
}

type childShape struct {
	Shape
	index    int
	localPos vec3.T
	localRot quat.T
	unsub    func()
}

var _ = checkShape(NewCompound())

func NewCompound() *Compound {
	cmp := object.NewComponent(&Compound{
		kind: CompoundShape,
	})
	cmp.Collider = newCollider(cmp, false)
	return cmp
}

func (c *Compound) colliderCreate() shapeHandle {
	return shape_new_compound(unsafe.Pointer(c))
}

func (c *Compound) colliderDestroy() {
	for _, shape := range c.shapes {
		shape.unsub()
	}
	c.shapes = nil
}

func (c *Compound) OnEnable() {
	c.Collider.OnEnable()
	c.colliderRefresh()
}

func (c *Compound) colliderRefresh() {
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
			localPos: shape.Transform().Position().Mul(c.Transform().WorldScale()),
			localRot: shape.Transform().Rotation(),
		}
		c.shapes = append(c.shapes, child)

		compound_add_child(c.handle, shape.shape(), child.localPos, child.localRot)

		// child shape changes should trigger a complete recreation of the compound shape
		unsubShape := shape.OnChange().Subscribe(func(s Shape) {
			c.refresh()
		})

		// adjust scale & local position on transform changes
		unsubTf := shape.Transform().OnChange().Subscribe(func(t transform.T) {
			scale := c.Transform().WorldScale()
			newPos := t.Position().Mul(scale)
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
}

func (c *Compound) OnDisable() {
	for _, shape := range c.shapes {
		shape.unsub()
	}
	c.shapes = nil
	c.Collider.OnDisable()
}

func (c *Compound) Name() string {
	return "CompoundShape"
}
