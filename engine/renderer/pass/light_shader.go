package pass

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/engine/renderer/uniform"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec4"
	"github.com/johanhenriksson/goworld/render/color"
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
	Depth    *descriptor.InputAttachment
	Camera   *descriptor.Uniform[uniform.Camera]
	Lights   *descriptor.Storage[uniform.Light]
	Shadow   *descriptor.SamplerArray
}

type LightConst struct {
	ViewProj    mat4.T
	Color       color.T
	Position    vec4.T
	Type        light.Type
	Index       uint32
	Range       float32
	Intensity   float32
	Attenuation light.Attenuation
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
	depthViews    []image.View
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
			Lights: &descriptor.Storage[uniform.Light]{
				Stages: core1_0.StageFragment,
				Size:   256,
			},
			Shadow: &descriptor.SamplerArray{
				Stages: core1_0.StageFragment,
				Count:  16,
			},
		})

	lightsh := mat.InstantiateMany(app.Pool(), app.Frames())

	diffuseViews := make([]image.View, app.Frames())
	normalViews := make([]image.View, app.Frames())
	positionViews := make([]image.View, app.Frames())
	depthViews := make([]image.View, app.Frames())
	for i := 0; i < app.Frames(); i++ {
		diffuseViews[i], _ = gbuffer.Diffuse()[i].View(gbuffer.Diffuse()[i].Format(), core1_0.ImageAspectColor)
		normalViews[i], _ = gbuffer.Normal()[i].View(gbuffer.Normal()[i].Format(), core1_0.ImageAspectColor)
		positionViews[i], _ = gbuffer.Position()[i].View(gbuffer.Position()[i].Format(), core1_0.ImageAspectColor)
		depthViews[i], _ = target.Depth()[i].View(target.Depth()[i].Format(), core1_0.ImageAspectDepth)

		lightDesc := lightsh[i].Descriptors()
		lightDesc.Diffuse.Set(diffuseViews[i])
		lightDesc.Normal.Set(normalViews[i])
		lightDesc.Position.Set(positionViews[i])
		lightDesc.Depth.Set(depthViews[i])
	}

	return &lightShader{
		mat:       mat,
		instances: lightsh,

		diffuseViews:  diffuseViews,
		normalViews:   normalViews,
		positionViews: positionViews,
		depthViews:    depthViews,
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
	for _, view := range ls.depthViews {
		view.Destroy()
	}
}
