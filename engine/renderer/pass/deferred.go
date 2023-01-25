package pass

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/sync"
	"github.com/johanhenriksson/goworld/render/vulkan"
)

type Deferred interface {
	Pass
	GBuffer() GeometryBuffer
}

type deferredPass struct {
	geometry *GeometryPass
	shadows  ShadowPass
}

func NewDeferredPass(target vulkan.Target, pool descriptor.Pool, geometrySubpass, shadowSubpass []DeferredSubpass) Deferred {
	shadows := NewShadowPass(target, pool, shadowSubpass)
	geometry := NewGeometryPass(target, pool, shadows, geometrySubpass)
	return &deferredPass{
		shadows:  shadows,
		geometry: geometry,
	}
}

func (d *deferredPass) Name() string {
	return "Deferred"
}

func (d *deferredPass) Completed() sync.Semaphore {
	return d.geometry.Completed()
}

func (d *deferredPass) Destroy() {
	d.shadows.Destroy()
	d.geometry.Destroy()
}

func (d *deferredPass) Draw(args render.Args, scene object.T) {
	d.shadows.Draw(args, scene)
	d.geometry.Draw(args, scene)
}

func (d *deferredPass) GBuffer() GeometryBuffer {
	return d.geometry.gbuffer
}
