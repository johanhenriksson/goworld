package physics

import (
	"log"
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
	cmp.Collider = newCollider(cmp)
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

func (c *Compound) colliderIsCompound() bool {
	return object.Get[*RigidBody](c) == nil && hasParentShape(c.Parent().Parent())
}

func (cmp *Compound) OnEnable() {
	cmp.Collider.OnEnable()

	// find all shapes that should be combined into the compound mesh
	// todo: also subscribe to attach/detach events on all relevant child objects
	shapes := Shapes(cmp)
	cmp.shapes = cmp.shapes[:0]
	for i, shape := range shapes {
		if shape.shape() == nil {
			panic("child shape is nil")
		}

		// find rotation relative to compound shape
		child := &childShape{
			Shape:    shape,
			index:    i,
			localPos: shape.Transform().Position().Mul(cmp.scale()),
			localRot: shape.Transform().Rotation(),
		}
		cmp.shapes = append(cmp.shapes, child)

		// shape_scaling_set(shape.shape(), localScale)
		compound_add_child(cmp.handle, shape.shape(), child.localPos, child.localRot)
		log.Println("add compound child", shape.Parent().Name(), child.localPos)

		// child shape changes should trigger a complete recreation of the compound shape
		unsubShape := shape.OnChange().Subscribe(func(s Shape) {
			cmp.refresh()
		})

		// adjust scale & local position on transform changes
		unsubTf := shape.Transform().OnChange().Subscribe(func(t transform.T) {
			newPos := t.Position().Mul(cmp.scale())
			newRot := t.Rotation()

			if !newPos.ApproxEqual(child.localPos) || !newRot.ApproxEqual(child.localRot) {
				compound_update_child(cmp.handle, child.index, newPos, newRot)
				child.localPos = newPos
				child.localRot = newRot
				cmp.OnChange().Emit(cmp)
			}
		})

		child.unsub = func() {
			unsubShape()
			unsubTf()
		}
	}
}

func (c *Compound) Name() string {
	return "CompoundShape"
}
