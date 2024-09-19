package pipeline

import (
	"reflect"

	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type Args struct {
	Pass      *renderpass.Renderpass
	Subpass   renderpass.Name
	Layout    *descriptor.SetLayout
	Shader    *shader.Shader
	Pointers  vertex.Pointers
	Constants []PushConstant

	Primitive       vertex.Primitive
	PolygonFillMode core1_0.PolygonMode
	CullMode        vertex.CullMode

	DepthTest  bool
	DepthWrite bool
	DepthClamp bool
	DepthFunc  core1_0.CompareOp

	StencilTest bool
}

func (args *Args) defaults() {
	if args.DepthFunc == 0 {
		args.DepthFunc = core1_0.CompareOpLessOrEqual
	}
	if args.Primitive == 0 {
		args.Primitive = vertex.Triangles
	}
}

type PushConstant struct {
	Stages core1_0.ShaderStageFlags
	Type   any
}

func (p *PushConstant) Size() int {
	t := reflect.TypeOf(p.Type)
	return int(t.Size())
}
