package pass

import (
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/pipeline"
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

type LightShader struct {
	pipe        *pipeline.Pipeline
	layout      *descriptor.Layout[*LightDescriptors]
	descriptors []*LightDescriptors

	diffuseTex   texture.Array
	normalTex    texture.Array
	positionTex  texture.Array
	occlusionTex texture.Array
}

func NewLightShader(app engine.App, pass *renderpass.Renderpass, gbuffer GeometryBuffer, occlusion engine.Target) *LightShader {
	dlayout := descriptor.NewLayout(app.Device(), "Lighting", &LightDescriptors{
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
	mat := pipeline.New(
		app.Device(),
		pipeline.Args{
			Shader:    app.Shaders().Fetch(shader.Ref("light")),
			Pass:      pass,
			Subpass:   LightingSubpass,
			Pointers:  vertex.ParsePointers(vertex.T{}),
			DepthTest: false,
		},
		dlayout)

	frames := gbuffer.Frames()

	var err error
	diffuseTex := make(texture.Array, frames)
	normalTex := make(texture.Array, frames)
	positionTex := make(texture.Array, frames)
	occlusionTex := make(texture.Array, frames)
	descriptors := make([]*LightDescriptors, frames)
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

		lightDesc := dlayout.Instantiate(app.Pool())
		lightDesc.Diffuse.Set(diffuseTex[i])
		lightDesc.Normal.Set(normalTex[i])
		lightDesc.Position.Set(positionTex[i])
		lightDesc.Occlusion.Set(occlusionTex[i])

		descriptors[i] = lightDesc
	}

	return &LightShader{
		pipe:        mat,
		layout:      dlayout,
		descriptors: descriptors,

		diffuseTex:   diffuseTex,
		normalTex:    normalTex,
		positionTex:  positionTex,
		occlusionTex: occlusionTex,
	}
}

func (ls *LightShader) Bind(cmd *command.Buffer, frame int) {
	cmd.CmdBindGraphicsPipeline(ls.pipe)
	cmd.CmdBindGraphicsDescriptor(0, ls.descriptors[frame])
}
func (ls *LightShader) Descriptors(frame int) *LightDescriptors {
	return ls.descriptors[frame]
}

func (ls *LightShader) Destroy() {
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
	for _, desc := range ls.descriptors {
		desc.Destroy()
	}
	ls.pipe.Destroy()
	ls.layout.Destroy()
}
