package pass

import (
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/pipeline"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vulkan"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type LightDescriptors struct {
	descriptor.Set
	Diffuse  *descriptor.InputAttachment
	Normal   *descriptor.InputAttachment
	Position *descriptor.InputAttachment
	Camera   *descriptor.Uniform[uniform.Camera]
	Lights   *descriptor.Storage[uniform.Light]
	Shadow   *descriptor.SamplerArray
}

type LightConst struct {
	Count uint32
}

type LightShader interface {
	Bind(command.Buffer, int)
	Descriptors(int) *LightDescriptors
	Destroy()
}

type lightShader struct {
	mat       material.T[*LightDescriptors]
	instances []material.Instance[*LightDescriptors]

	diffuseViews  []image.View
	normalViews   []image.View
	positionViews []image.View
}

func NewLightShader(app vulkan.App, pass renderpass.T, gbuffer GeometryBuffer) LightShader {
	mat := material.New(
		app.Device(),
		material.Args{
			Shader:   app.Shaders().Fetch(shader.NewRef("light")),
			Pass:     pass,
			Subpass:  LightingSubpass,
			Pointers: vertex.ParsePointers(vertex.T{}),
			Constants: []pipeline.PushConstant{
				{
					Stages: core1_0.StageFragment,
					Type:   LightConst{},
				},
			},
			DepthTest: true,
		},
		&LightDescriptors{
			Diffuse: &descriptor.InputAttachment{
				Stages: core1_0.StageFragment,
			},
			Normal: &descriptor.InputAttachment{
				Stages: core1_0.StageFragment,
			},
			Position: &descriptor.InputAttachment{
				Stages: core1_0.StageFragment,
			},
			Camera: &descriptor.Uniform[uniform.Camera]{
				Stages: core1_0.StageFragment,
			},
			Lights: &descriptor.Storage[uniform.Light]{
				Stages: core1_0.StageFragment,
				Size:   256,
			},
			Shadow: &descriptor.SamplerArray{
				Stages: core1_0.StageFragment,
				Count:  32,
			},
		})

	frames := gbuffer.Frames()
	lightsh := mat.InstantiateMany(app.Pool(), frames)

	diffuseViews := make([]image.View, frames)
	normalViews := make([]image.View, frames)
	positionViews := make([]image.View, frames)
	for i := 0; i < frames; i++ {
		diffuseViews[i], _ = gbuffer.Diffuse()[i].View(gbuffer.Diffuse()[i].Format(), core1_0.ImageAspectColor)
		normalViews[i], _ = gbuffer.Normal()[i].View(gbuffer.Normal()[i].Format(), core1_0.ImageAspectColor)
		positionViews[i], _ = gbuffer.Position()[i].View(gbuffer.Position()[i].Format(), core1_0.ImageAspectColor)

		lightDesc := lightsh[i].Descriptors()
		lightDesc.Diffuse.Set(diffuseViews[i])
		lightDesc.Normal.Set(normalViews[i])
		lightDesc.Position.Set(positionViews[i])
	}

	return &lightShader{
		mat:       mat,
		instances: lightsh,

		diffuseViews:  diffuseViews,
		normalViews:   normalViews,
		positionViews: positionViews,
	}
}

func (ls *lightShader) Bind(buf command.Buffer, frame int) {
	ls.instances[frame].Bind(buf)
}
func (ls *lightShader) Descriptors(frame int) *LightDescriptors {
	return ls.instances[frame].Descriptors()
}

func (ls *lightShader) Destroy() {
	ls.mat.Destroy()
	for _, view := range ls.diffuseViews {
		view.Destroy()
	}
	for _, view := range ls.normalViews {
		view.Destroy()
	}
	for _, view := range ls.positionViews {
		view.Destroy()
	}
}
