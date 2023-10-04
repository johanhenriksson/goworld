package player

import (
	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/geometry/sphere"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type Skydome struct {
	object.Object
	*sphere.Mesh
}

func NewSkydome() *Skydome {
	dome := object.New("Skydome", &Skydome{
		Mesh: sphere.New(&material.Def{
			Pass:         material.Forward,
			Shader:       "forward/skybox",
			VertexFormat: vertex.T{},
			DepthTest:    true,
			DepthWrite:   true,
			DepthFunc:    core1_0.CompareOpLessOrEqual,
			Primitive:    vertex.Triangles,
			CullMode:     vertex.CullFront,
			Transparent:  true,
		}),
	})
	dome.SetShadows(false)
	return dome
}

func (d *Skydome) EditorUpdate(scene object.Component, dt float32) {
	d.Update(scene, dt)
}

func (d *Skydome) Update(scene object.Component, dt float32) {
	d.Object.Update(scene, dt)

	cam := object.GetInParents[*camera.Camera](d)
	if cam == nil {
		return
	}

	d.Transform().SetScale(vec3.New(cam.Far, cam.Far, cam.Far))
}
