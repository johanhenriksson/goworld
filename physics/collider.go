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
	colliderImpl

	handle    shapeHandle
	changed   *events.Event[Shape]
	scaled    bool
	lastScale vec3.T

	unsubTf func()
}

type colliderImpl interface {
	colliderCreate() shapeHandle
	colliderRefresh()
	colliderDestroy()
	colliderIsCompound() bool
}

var _ Shape = &Collider{}

func newCollider(impl colliderImpl, scaled bool) *Collider {
	col := object.NewComponent(&Collider{
		colliderImpl: impl,
		changed:      events.New[Shape](),
		scaled:       scaled,
	})

	// create initial handle
	col.handle = col.colliderCreate()

	runtime.SetFinalizer(col, func(b *Collider) {
		b.destroy()
	})
	return col
}

func (c *Collider) OnChange() *events.Event[Shape] {
	return c.changed
}

func (c *Collider) shape() shapeHandle {
	return c.handle
}

func (c *Collider) scale() vec3.T {
	if c.scaled {
		return c.Transform().WorldScale()
	}
	return vec3.One
}

func (c *Collider) OnEnable() {
	// subscribe to transform updates so that we may react to scale changes
	if c.unsubTf != nil {
		panic("should not be subscribed")
	}
	c.unsubTf = c.Transform().OnChange().Subscribe(c.transformRefresh)
	c.transformRefresh(c.Transform())
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
	c.OnChange().Emit(c)
}

func (c *Collider) refresh() {
	c.destroy()
	c.handle = c.colliderCreate()
	c.colliderRefresh()
	shape_scaling_set(c.handle, c.scale())
	c.OnChange().Emit(c)
}

func (c *Collider) destroy() {
	c.colliderDestroy()
	if c.handle != nil {
		shape_delete(&c.handle)
	}
}
