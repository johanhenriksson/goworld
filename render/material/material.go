package material

import (
	"fmt"
	"log"

	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/pipeline"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/vkngwrapper/core/v2/core1_0"
)

// Materials combine pipelines and descriptors into a common unit.
type Material[D descriptor.Set] struct {
	device  *device.Device
	shader  *shader.Shader
	playout *pipeline.Layout
	pipe    *pipeline.Pipeline
	pass    *renderpass.Renderpass
}

type Args struct {
	Shader    *shader.Shader
	Pass      *renderpass.Renderpass
	Subpass   renderpass.Name
	Constants []pipeline.PushConstant

	Pointers   vertex.Pointers
	Primitive  vertex.Primitive
	DepthTest  bool
	DepthWrite bool
	DepthClamp bool
	DepthBias  float32
	DepthSlope float32
	DepthFunc  core1_0.CompareOp
	CullMode   vertex.CullMode
}

func New[D descriptor.Set](device *device.Device, args Args, descriptors ...descriptor.SetLayout) *Material[D] {
	if device == nil {
		panic("device is nil")
	}
	if args.Shader == nil {
		panic("shader is nil")
	}

	for i, ptr := range args.Pointers {
		if index, kind, exists := args.Shader.Input(ptr.Name); exists {
			ptr.Bind(index, kind)
			args.Pointers[i] = ptr
		} else {
			log.Printf("no attribute in shader %s\n", ptr.Name)
		}
	}

	if args.Primitive == 0 {
		args.Primitive = vertex.Triangles
	}

	// create new descriptor set layout
	// ... this could be cached ...
	// todo: this ought to happen somewhere else
	// descLayout := descriptor.NewLayout(device, args.Shader.Name(), descriptors)

	// crete pipeline pipelineLayout
	// ... this could be cached ...
	pipelineLayout := pipeline.NewLayout(device, descriptors, args.Constants)

	pipelineName := fmt.Sprintf("%s/%s", args.Pass.Name(), args.Shader.Name())
	pipe := pipeline.New(device, pipeline.Args{
		Key:      pipelineName,
		Layout:   pipelineLayout,
		Pass:     args.Pass,
		Subpass:  args.Subpass,
		Shader:   args.Shader,
		Pointers: args.Pointers,

		Primitive:  args.Primitive,
		DepthTest:  args.DepthTest,
		DepthWrite: args.DepthWrite,
		DepthClamp: args.DepthClamp,
		DepthFunc:  args.DepthFunc,
		CullMode:   args.CullMode,
	})

	return &Material[D]{
		device:  device,
		shader:  args.Shader,
		playout: pipelineLayout,
		pipe:    pipe,
		pass:    args.Pass,
	}
}

func (m *Material[D]) Bind(cmd *command.Buffer) {
	cmd.CmdBindGraphicsPipeline(m.pipe)
}

func (m *Material[D]) TextureSlots() []texture.Slot {
	return m.shader.Textures()
}

func (m *Material[D]) Destroy() {
	if m.playout != nil {
		m.playout.Destroy()
		m.playout = nil
	}
	if m.pipe != nil {
		m.pipe.Destroy()
		m.pipe = nil
	}
}
