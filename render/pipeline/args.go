package pipeline

import (
	"reflect"

	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/vertex"
	vk "github.com/vulkan-go/vulkan"
)

type Args struct {
	Pass     renderpass.T
	Subpass  string
	Layout   Layout
	Shader   shader.T
	Pointers vertex.Pointers

	Primitive       vertex.Primitive
	PolygonFillMode vk.PolygonMode
	CullMode        vk.CullModeFlagBits

	DepthTest  bool
	DepthWrite bool
	DepthFunc  vk.CompareOp

	StencilTest bool
}

func (args *Args) defaults() {
	if args.DepthFunc == 0 {
		args.DepthFunc = vk.CompareOpLessOrEqual
	}
	if args.Primitive == 0 {
		args.Primitive = vertex.Triangles
	}
}

type PushConstant struct {
	Stages vk.ShaderStageFlagBits
	Type   any
}

func (p *PushConstant) Size() int {
	t := reflect.TypeOf(p.Type)
	return int(t.Size())
}
