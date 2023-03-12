package pass

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/engine/renderer/uniform"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec4"
	"github.com/johanhenriksson/goworld/render/color"
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
	Depth    *descriptor.InputAttachment
	Camera   *descriptor.Uniform[uniform.Camera]
	Shadow   *descriptor.SamplerArray
}

type LightConst struct {
	ViewProj    mat4.T
	Color       color.T
	Position    vec4.T
	Type        light.Type
	Shadowmap   uint32
	Range       float32
	Intensity   float32
	Attenuation light.Attenuation
}

type LightShader interface {
	material.Instance[*LightDescriptors]
	Destroy()
}

type lightShader struct {
	material.Instance[*LightDescriptors]

	diffuseView  image.View
	normalView   image.View
	positionView image.View
	depthView    image.View
}

func NewLightShader(app vulkan.App, pass renderpass.T, target RenderTarget, gbuffer GeometryBuffer) LightShader {
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
			Depth: &descriptor.InputAttachment{
				Stages: core1_0.StageFragment,
			},
			Camera: &descriptor.Uniform[uniform.Camera]{
				Stages: core1_0.StageFragment,
			},
			Shadow: &descriptor.SamplerArray{
				Stages: core1_0.StageFragment,
				Count:  16,
			},
		})

	lightsh := mat.Instantiate(app.Pool())

	diffuseView, _ := gbuffer.Diffuse().View(gbuffer.Diffuse().Format(), core1_0.ImageAspectColor)
	normalView, _ := gbuffer.Normal().View(gbuffer.Normal().Format(), core1_0.ImageAspectColor)
	positionView, _ := gbuffer.Position().View(gbuffer.Position().Format(), core1_0.ImageAspectColor)
	depthView, _ := target.Depth()[0].View(target.Depth()[0].Format(), core1_0.ImageAspectDepth)

	lightDesc := lightsh.Descriptors()
	lightDesc.Diffuse.Set(diffuseView)
	lightDesc.Normal.Set(normalView)
	lightDesc.Position.Set(positionView)
	lightDesc.Depth.Set(depthView)

	return &lightShader{
		Instance: lightsh,

		diffuseView:  diffuseView,
		normalView:   normalView,
		positionView: positionView,
		depthView:    depthView,
	}
}

func (ls *lightShader) Destroy() {
	ls.Instance.Material().Destroy()
	ls.diffuseView.Destroy()
	ls.normalView.Destroy()
	ls.positionView.Destroy()
	ls.depthView.Destroy()
}
