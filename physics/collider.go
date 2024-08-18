package physics

import (
	"runtime"

	"github.com/johanhenriksson/goworld/core/events"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/transform"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Collider struct {
	object.Component

	cmp colliderImpl

	handle    shapeHandle
	changed   events.Event[Shape]
	scaled    bool
	lastScale vec3.T

	unsubTf func()
}

type colliderImpl interface {
	Shape

	colliderCreate() shapeHandle
	colliderRefresh()
	colliderDestroy()
	colliderIsCompound() bool
}

func newCollider(pool object.Pool, impl colliderImpl, scaled bool) *Collider {
	col := object.NewComponent(pool, &Collider{
		cmp:    impl,
		scaled: scaled,
	})

	// create initial handle
	col.handle = col.cmp.colliderCreate()

	runtime.SetFinalizer(col, func(b *Collider) {
		b.destroy()
	})
	return col
}

func (c *Collider) OnChange() *events.Event[Shape] {
	return &c.changed
}

func (c *Collider) shape() shapeHandle {
	return c.handle
}

func (c *Collider) scale() vec3.T {
	if c.scaled {
		return c.cmp.Transform().WorldScale()
	}
	return vec3.One
}

func (c *Collider) OnEnable() {
	// subscribe to transform updates so that we may react to scale changes
	if c.unsubTf != nil {
		panic("should not be subscribed")
	}
	c.unsubTf = c.cmp.Transform().OnChange().Subscribe(c.transformRefresh)
	c.transformRefresh(c.cmp.Transform())
}

func (c *Collider) OnDisable() {
	c.unsubTf()
	c.unsubTf = nil
}

func (c *Collider) transformRefresh(t transform.T) {
	newScale := c.scale()
	if newScale == c.lastScale {
		return
	}

	c.lastScale = newScale
	shape_scaling_set(c.handle, newScale)

	// raising OnChange is technically not required since we dont recreate the shape
	c.OnChange().Emit(c.cmp)
}

func (c *Collider) refresh() {
	c.destroy()
	c.handle = c.cmp.colliderCreate()
	c.cmp.colliderRefresh()
	shape_scaling_set(c.handle, c.scale())
	c.OnChange().Emit(c.cmp)
}

func (c *Collider) destroy() {
	c.cmp.colliderDestroy()
	if c.handle != nil {
		shape_delete(&c.handle)
	}
}
