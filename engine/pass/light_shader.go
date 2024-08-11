package pass

import (
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type LightDescriptors struct {
	descriptor.Set
	Camera    *descriptor.Uniform[uniform.Camera]
	Lights    *descriptor.Storage[uniform.Light]
	Diffuse   *descriptor.Sampler
	Normal    *descriptor.Sampler
	Position  *descriptor.Sampler
	Occlusion *descriptor.Sampler
	Shadow    *descriptor.SamplerArray
}

type LightShader interface {
	Bind(command.Buffer, int)
	Descriptors(int) *LightDescriptors
	Destroy()
}

type lightShader struct {
	mat       *material.Material[*LightDescriptors]
	instances []*material.Instance[*LightDescriptors]

	diffuseTex   []texture.T
	normalTex    []texture.T
	positionTex  []texture.T
	occlusionTex []texture.T
}

func NewLightShader(app engine.App, pass renderpass.T, gbuffer GeometryBuffer, occlusion engine.Target) LightShader {
	mat := material.New(
		app.Device(),
		material.Args{
			Shader:    app.Shaders().Fetch(shader.NewRef("light")),
			Pass:      pass,
			Subpass:   LightingSubpass,
			Pointers:  vertex.ParsePointers(vertex.T{}),
			DepthTest: false,
		},
		&LightDescriptors{
			Camera: &descriptor.Uniform[uniform.Camera]{
				Stages: core1_0.StageFragment,
			},
			Lights: &descriptor.Storage[uniform.Light]{
				Stages: core1_0.StageFragment,
				Size:   256,
			},
			Diffuse: &descriptor.Sampler{
				Stages: core1_0.StageFragment,
			},
			Normal: &descriptor.Sampler{
				Stages: core1_0.StageFragment,
			},
			Position: &descriptor.Sampler{
				Stages: core1_0.StageFragment,
			},
			Occlusion: &descriptor.Sampler{
				Stages: core1_0.StageFragment,
			},
			Shadow: &descriptor.SamplerArray{
				Stages: core1_0.StageFragment,
				Count:  32,
			},
		})

	frames := gbuffer.Frames()
	lightsh := mat.InstantiateMany(app.Pool(), frames)

	var err error
	diffuseTex := make([]texture.T, frames)
	normalTex := make([]texture.T, frames)
	positionTex := make([]texture.T, frames)
	occlusionTex := make([]texture.T, frames)
	for i := 0; i < frames; i++ {
		diffuseTex[i], err = texture.FromImage(app.Device(), "deferred-diffuse", gbuffer.Diffuse()[i], texture.Args{
			Filter: texture.FilterNearest,
		})
		if err != nil {
			panic(err)
		}
		normalTex[i], err = texture.FromImage(app.Device(), "deferred-normal", gbuffer.Normal()[i], texture.Args{
			Filter: texture.FilterNearest,
		})
		if err != nil {
			panic(err)
		}
		positionTex[i], err = texture.FromImage(app.Device(), "deferred-position", gbuffer.Position()[i], texture.Args{
			Filter: texture.FilterNearest,
		})
		if err != nil {
			panic(err)
		}
		occlusionTex[i], err = texture.FromImage(app.Device(), "deferred-ssao", occlusion.Surfaces()[i], texture.Args{
			Filter: texture.FilterNearest,
		})
		if err != nil {
			panic(err)
		}

		lightDesc := lightsh[i].Descriptors()
		lightDesc.Diffuse.Set(diffuseTex[i])
		lightDesc.Normal.Set(normalTex[i])
		lightDesc.Position.Set(positionTex[i])
		lightDesc.Occlusion.Set(occlusionTex[i])
	}

	return &lightShader{
		mat:       mat,
		instances: lightsh,

		diffuseTex:   diffuseTex,
		normalTex:    normalTex,
		positionTex:  positionTex,
		occlusionTex: occlusionTex,
	}
}

func (ls *lightShader) Bind(buf command.Buffer, frame int) {
	ls.instances[frame].Bind(buf)
}
func (ls *lightShader) Descriptors(frame int) *LightDescriptors {
	return ls.instances[frame].Descriptors()
}

func (ls *lightShader) Destroy() {
	for _, view := range ls.diffuseTex {
		view.Destroy()
	}
	for _, view := range ls.normalTex {
		view.Destroy()
	}
	for _, view := range ls.positionTex {
		view.Destroy()
	}
	for _, view := range ls.occlusionTex {
		view.Destroy()
	}
	ls.mat.Destroy()
}
