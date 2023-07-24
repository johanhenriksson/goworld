package physics

import (
	"log"
	"runtime"

	"github.com/johanhenriksson/goworld/core/events"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/transform"
	"github.com/johanhenriksson/goworld/math/vec3"
)

func defaultCompoundCheck(c object.Component) bool {
	return hasParentShape(c.Parent()) || hasParentShape(c.Parent().Parent())
}

type Collider struct {
	object.Component
	colliderImpl

	handle    shapeHandle
	changed   *events.Event[Shape]
	compound  bool
	lastScale vec3.T

	unsubTf func()
}

type colliderImpl interface {
	colliderCreate() shapeHandle
	colliderDestroy()
	colliderIsCompound() bool
}

var _ Shape = &Collider{}

func newCollider(impl colliderImpl) *Collider {
	col := object.NewComponent(&Collider{
		colliderImpl: impl,
		changed:      events.New[Shape](),
		lastScale:    vec3.One,
	})

	// trigger initial resize
	col.refresh()

	runtime.SetFinalizer(col, func(b *Collider) {
		b.destroy()
	})
	return col
}

func (s *Collider) OnChange() *events.Event[Shape] {
	return s.changed
}

func (s *Collider) shape() shapeHandle {
	return s.handle
}

func (b *Collider) scale() vec3.T {
	if b.compound {
		return b.Transform().Scale()
	}
	return b.Transform().WorldScale()
}

func (c *Collider) OnEnable() {
	// check if the collider is part of a compound shape.
	// if it is, it should be scaled according to its local scale factor
	c.compound = c.colliderIsCompound()
	log.Println("collider", c.Parent().Name(), "is compound:", c.compound)

	c.unsubTf = c.Transform().OnChange().Subscribe(c.transformRefresh)
	c.transformRefresh(c.Transform())
}

func (b *Collider) OnDisable() {
	b.unsubTf()
}

func (c *Collider) transformRefresh(t transform.T) {
	newScale := c.scale()
	if newScale == c.lastScale {
		return
	}

	log.Println("collider scale update", c.Parent().Name(), ":", newScale)
	c.lastScale = newScale
	shape_scaling_set(c.handle, newScale)

	// raising OnChange is technically not required since we dont recreate the shape
	c.OnChange().Emit(c)
}

func (b *Collider) refresh() {
	b.destroy()
	b.handle = b.colliderCreate()
	shape_scaling_set(b.handle, b.scale())
	b.OnChange().Emit(b)
}

func (b *Collider) destroy() {
	b.colliderDestroy()
	if b.handle != nil {
		shape_delete(&b.handle)
	}
}
